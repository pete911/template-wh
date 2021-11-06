## Usage

[Helm](https://helm.sh) must be installed to use the charts.  Please refer to
Helm's [documentation](https://helm.sh/docs) to get started.

Once Helm has been set up correctly, add the repo as follows:

  helm repo add pete911-template-wh https://pete911.github.io/template-wh

If you had already added this repo earlier, run `helm repo update` to retrieve
the latest versions of the packages.  You can then run `helm search repo
pete911-template-wh` to see the charts.

To install the template-wh chart:

    helm install template-wh pete911-template-wh/template-wh

To uninstall the chart:

    helm delete template-wh

