![Verdex screenshot](./docs/images/screenshot.png)

**Verdex** is a fast, collaborative, and open-source tool designed to 
**detect online service versions through feature-based analysis**,
even when traditional versioning information is unavailable.

🔧 Instead of relying on exposed headers, Verdex quickly analyzes the features present on a service. 
For example, **if a feature was introduced in a specific version, Verdex uses this information to deduce the 
minimum version of the service** (see [How It Works](https://verdex-docs.mintlify.app/how-it-works) for more information).

🪶 Verdex relies solely on **non-aggressive methods**, in particular by using mostly static assets (javascript, css and html files).

🎯 Currently, Verdex supports detecting versions of **Keycloak** and **NextCloud** (see [Products](https://verdex-docs.mintlify.app/essentials/products) documentation).

## 📖 Documentation
The entire tool documentation is available at [https://verdex-docs.mintlify.app](https://verdex-docs.mintlify.app)

## 🚀 Install Verdex

Using Go:
```bash
go install -v github.com/verdexlab/verdex@latest
```

or using Docker:
```bash
echo 'alias verdex="docker run --rm verdexlab/verdex:latest"' >> ~/.bashrc
source ~/.bashrc
```

> 📖 Read documentation for more installation options: 
[https://verdex-docs.mintlify.app/quickstart](https://verdex-docs.mintlify.app/quickstart)

## 🌪️ Run version detection

```bash
verdex -target https://keycloak.example.com -output-json ./results.json
```

> 📖 Read documentation for more detection options: 
[https://verdex-docs.mintlify.app/essentials/options](https://verdex-docs.mintlify.app/essentials/options)

## License
Verdex is distributed under MIT License.
