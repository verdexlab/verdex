import requests
from urllib.parse import urlencode

github_organization = "verdexlab"
github_repository = "verdex"
github_issues_label = "product-suggestion"
documentation_snippet_path = "docs/snippets/products-votes.mdx"


# Fetch products votes issues from GitHub
query_params = {
    "state": "open",
    "per_page": 100,
    "labels": github_issues_label,
}

res = requests.get(f"https://api.github.com/repos/{github_organization}/{github_repository}/issues?{urlencode(query_params)}")
if res.status_code != 200:
    print("status_code", res.status_code)
    print("response", res.text)
    raise Exception("Failed to fetch products votes issues from GitHub")
issues = res.json()
print(f"{len(issues)} issues found with label #{github_issues_label} on {github_organization}/{github_repository}")


# Generate documentation cards
docs_contents = []
for issue in issues:
    product_name = ""
    product_description = ""
    product_icon = ""

    for body_part in issue["body"].split("\n\n### "):
        if body_part.startswith('### Product Name\n\n'):
            product_name = body_part.split('### Product Name\n\n')[1]
        elif body_part.startswith('Product Description\n\n'):
            product_description = body_part.split('Product Description\n\n')[1]
        elif body_part.startswith('Product Icon\n\n'):
            product_icon = body_part.split('Product Icon\n\n')[1]

    if not product_name or not product_description or not product_icon:
        print(f"Skipped issue #{issue['number']} (failed to parse product name, description and icon)")
        continue

    votes = 0
    product_reactions = []
    if issue["reactions"]["+1"] > 0:
        product_reactions.append(f"👍 {issue['reactions']['+1']}")
        votes += issue["reactions"]["+1"]
    if issue["reactions"]["hooray"] > 0:
        product_reactions.append(f"🎉 {issue['reactions']['hooray']}")
        votes += issue["reactions"]["hooray"]
    if issue["reactions"]["heart"] > 0:
        product_reactions.append(f"❤️ {issue['reactions']['heart']}")
        votes += issue["reactions"]["heart"]
    if issue["reactions"]["rocket"] > 0:
        product_reactions.append(f"🚀 {issue['reactions']['rocket']}")
        votes += issue["reactions"]["rocket"]
    if len(product_reactions) == 0:
        product_reactions.append("No votes")

    product_description += f"<br/><br/><a href=\"{issue['html_url']}\" target=\"_blank\"><Tooltip tip=\"Click to vote for this product\">{'&nbsp;&nbsp;&nbsp;'.join(product_reactions)}</Tooltip></a>"
    docs_contents.append({
        "content": f"""<Card title="{product_name}" icon="{product_icon}" iconType="duotone">\n  {product_description}\n</Card>\n""",
        "votes": votes,
    })

    print(f"Added issue #{issue['number']} ({votes} votes)")

# Sort by votes DESC (most voted first)
docs_contents = sorted(docs_contents, key=lambda d: d['votes'])
docs_contents.reverse()

f = open(documentation_snippet_path, "w")
f.write("".join([c['content'] for c in docs_contents]))
f.close()
