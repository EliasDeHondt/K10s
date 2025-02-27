![logo](/Assets/logo.png)
# 🤍🩵💜README💜🩵🤍

## 📘Table of Contents

1. [📘Table of Contents](#📘table-of-contents)
2. [🖖Introduction](#🖖introduction)
    1. [👉Key Features](#👉key-features)
3. [🎨Color Scheme](#🎨color-scheme)
4. [🎥Video](#🎥video)
5. [✒️Pitch](#✒️pitch)
6. [📷Photos](#📷photos)
7. [📚How to deploy and destroy](#📚how-to-deploy-and-destroy)
    1. [🌌Kubernetes](#🌌kubernetes)
    2. [🐳Docker](#🐳docker)

---

## 🖖Introduction

**K10s** is an open-source project designed to provide a comprehensive visual and hierarchical representation of Kubernetes clusters, including superclusters (clusters of clusters). This tool aims to give users an intuitive understanding of how their clusters operate, manage resources, and interconnect with underlying components.

> If you find this project helpful or interesting, feel free to ⭐️ **star this repository** on GitHub. Your support helps us grow and reach more users in the community!

### 👉Key Features
- Visualize the full hierarchy of your Kubernetes ecosystem, from superclusters to individual pools, nodes and pods...
- Gain real-time insights into cluster performance and management.
- Simplify complex Kubernetes structures into an understandable format for developers and administrators alike.

To ensure proper use and contribution, please refer to the following documentation:
- 📜 [LICENSE](/LICENSE.md): Details on permissions and limitations for using the software.
- 🔒 [SECURITY](/SECURITY.md): Guidelines for identifying and reporting vulnerabilities.
- 🤝 [CONTRIBUTING](/CONTRIBUTING.md): Instructions for contributing to the project and submitting pull requests.
- 📝 [CODE OF CONDUCT](/CODE-OF-CONDUCT.md): Guidelines for participating in the community and maintaining a respectful environment.

## 🎥Video

[Watch our video on the K10s website](https://k10s.eliasdh.com)

## ✒️Pitch

- 💜K10s is your go-to tool for effortlessly exploring the hierarchy and structure of a Kubernetes clusters.🤍 By harnessing the power of the Kubernetes API, K10s delivers an intuitive and detailed overview of your cluster architecture.🩵

## 📷Photos

![photo 1](/Assets/photo1.png)

![photo 2](/Assets/photo2.png)

## 🎨Color Scheme

The project's visual identity is defined by the following colors:

```css
:root[data-theme="dark"] {
    --text: #ecf1f4;
    --background: #040c10;
    --primary: #004b7a;
    --secondary: #3c3154;
    --tertiary: #262b2e;
    --quaternary: #1f2326;
    --accent: #5f4270;
}

:root[data-theme="light"] {
    --text: #0b1013;
    --background: #eff7fb;
    --primary: #85d0ff;
    --secondary: #b6abce;
    --tertiary: #e8eff3;
    --quaternary: #dedfed;
    --accent: #ac8fbd;
}
```

## 📚How to deploy and destroy

### 🌌Kubernetes

- Step 1 - Deploy the application:
```bash
kubectl apply -f https://raw.githubusercontent.com/EliasDeHondt/K10s/refs/heads/main/Kubernetes/k10s.yaml
```

- Step 2 - Create the necessary secrets for the application (replace the values with your own):
```bash
kubectl create secret generic k10s-secret-user -n k10s-namespaces --from-literal=USERNAME=admin --from-literal=PASSWORD=admin
kubectl create secret generic k10s-secret-jwt -n k10s-namespaces --from-literal=JWT_SECRET=$(openssl rand -base64 32) # apt install openssl.
```

- Step 3 - Get the external IP of the application:
```bash
kubectl get svc -n k10s-namespaces
```

---

- If you want to delete the application:
```bash
kubectl delete -f https://raw.githubusercontent.com/EliasDeHondt/K10s/refs/heads/main/Kubernetes/k10s.yaml
```

### 🐳Docker

> **NOTE:** The following commands are simply for testing purposes only. For production use, please refer to the Kubernetes deployment instructions above.

- Pull the latest image and run the container
```bash
sudo docker pull ghcr.io/eliasdehondt/k10s-frontend:latest
sudo docker run --name k10s-frontend-container -p 80:80 -d ghcr.io/eliasdehondt/k10s-frontend:latest
```

---

- Stop and remove the existing container and image
```bash
sudo docker stop k10s-frontend-container
sudo docker rm k10s-frontend-container
sudo docker rmi ghcr.io/eliasdehondt/k10s-frontend:latest
```


## 📚How to run K10s locally for development purposes

- Step 1 - Clone the repository:
```bash
git clone https://github.com/EliasDeHondt/K10s.git
```

- Step 2 - Install the necessary dependencies for the frontend:
```bash
cd /K10s/App/Frontend
npm install
```

- Step 3 - Run the frontend:
```bash
ng serve
```

- Step 4 - Install the necessary dependencies for the backend:
```bash
cd /K10s/App/Backend/cmd
go build main.go
```

- Step 5 - Run the backend:
```bash
go run main.go
```