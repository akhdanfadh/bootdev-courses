# Learn Kubernetes

This course is a bit different than the other courses of Boot.dev. We'll be doing *very little* coding in the browser. [Kubernetes](https://kubernetes.io/) is a distributed system of servers that host software applications, and you interact with it primarily through your local command line - it's not a programming language.

As such, you'll complete the majority of this course on your own machine. We'll use a combination of HTTP-based tests and quizzes to ensure you're on the right track.

## What We'll Need

- The Kubernetes command-line tool, [`kubectl`](https://kubernetes.io/docs/tasks/tools/): allows us to run commands against Kubernetes clusters. It's a client that communicates with a Kubernetes API server. For Homebrew users, simply run `brew install kubectl`.
- [Minikube](https://minikube.sigs.k8s.io/docs/), a fantastic tool that allows us to run a single-node Kubernetes cluster on our local machine. For Homebrew users, run `brew install minikube`, then start with `minikube start --extra-config "apiserver.cors-allowed-origins=["http://boot.dev"]"` to allow boot.dev to access our local cluster.
