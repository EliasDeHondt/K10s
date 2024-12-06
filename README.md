![logo](/Images/icon-192x192.png)
# 🤍🩵💜README💜🩵🤍

## 📘Table of Contents

1. [📘Table of Contents](#📘table-of-contents)
2. [🖖Introduction](#🖖introduction)
    1. [👉Key Features:](#👉key-features)
3. [🎨Color Scheme](#🎨color-scheme)
4. [📚How to Deploy](#📚how-to-deploy)

---

## 🖖Introduction

**K10s** is an open-source project designed to provide a comprehensive visual and hierarchical representation of Kubernetes clusters, including superclusters (clusters of clusters). This tool aims to give users an intuitive understanding of how their clusters operate, manage resources, and interconnect with underlying components.

> If you find this project helpful or interesting, feel free to ⭐️ **star this repository** on GitHub. Your support helps us grow and reach more users in the community!

### 👉Key Features:  
- Visualize the full hierarchy of your Kubernetes ecosystem, from superclusters to individual pools, nodes and pods...
- Gain real-time insights into cluster performance and management.
- Simplify complex Kubernetes structures into an understandable format for developers and administrators alike.

To ensure proper use and contribution, please refer to the following documentation:
- 📜 [LICENSE](/LICENSE.md): Details on permissions and limitations for using the software.
- 🔒 [SECURITY](/SECURITY.md): Guidelines for identifying and reporting vulnerabilities.
- 🤝 [CONTRIBUTING](/CONTRIBUTING.md): Instructions for contributing to the project and submitting pull requests.
- 📝 [CODE OF CONDUCT](/CODE-OF-CONDUCT.md): Guidelines for participating in the community and maintaining a respectful environment.

## 🎨Color Scheme

The project's visual identity is defined by the following colors:

- **Primary Color**: <span style="color: #85d0ff;">(#85d0ff)</span>
- **Secondary Color**: <span style="color: #b6abce;">(#b6abce)</span>
- **Accent Color**: <span style="color: #ac8fbd;">(#ac8fbd)</span>
- **Background Color**: <span style="color: #eff7fb;">(#eff7fb)</span>
- **Text Color**: <span style="color: #0b1013;">(#0b1013)</span>

## 📚How to deploy

- Step 1: Deploy the application:
```bash
kubectl apply -f https://raw.githubusercontent.com/EliasDeHondt/K10s/refs/heads/main/Kubernetes/k10s.yaml
```

- Step 2: Get the external IP of the application:
```bash
kubectl get svc -n k10s
```

- If you want to delete the application:
```bash
kubectl delete -f https://raw.githubusercontent.com/EliasDeHondt/K10s/refs/heads/main/Kubernetes/k10s.yaml
```