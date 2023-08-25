# Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit <https://cla.opensource.microsoft.com>.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Types of Contributions

We welcome all types of contributions.  Bug fixes and documentation updates are great ways to get started contributing. If you're looking to make a more substantial contribution to the provider, consider one of the following options.

### Resources

Creating a new [resource](https://developer.hashicorp.com/terraform/plugin/framework/resources) can allow terraform to manage new Power Platform infrastructure not currently provided by the provider.

### Data Sources

Creating a new [data source](https://developer.hashicorp.com/terraform/plugin/framework/data-sources) can allow terraform to reference data about Power Platform services and infrastrucutre.

### Guides

Guides provide documentation and example Terraform files showing how to use Terraform and the Power Platform Terraform Provider to implement more complex scenarios.

## Pull Request Checklist

PRs for new resources or data sources are expected to meeting the following criteria:

- Add a production quality implementation of the resource or datasource in [/internal](/internal/)
- Add unit tests and acceptance tests for your contribution in [/internal](/internal/).  
    - Tests should pass and provide >90% coverage of your contribution
- Add examples for your contribution in [/examples](/examples/) (see [Terraform Documentaion on examples](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation#add-configuration-examples))
- Add [schema descriptions](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation#add-schema-descriptions) for your resource or data source in [/internal](/internal/) 
- and/or [/templates](/templates/)
- [Update auto-generated documentation](#updating-documentation) in [/docs](/docs/). (Don't manually edit [/docs](/docs/) or your updates will be overwritten)
- Ensure the PR description clearly describes the feature you're adding and any known limitations
