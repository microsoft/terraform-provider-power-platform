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

We welcome all types of contributions.  Bug fixes and documentation updates are great ways to get started contributing. If you're looking to make a more substantial contribution to the provider, consider one of the following options (in order of difficulty).

### Examples

Contributing examples of using the Power Platform Terraform Provider is a great way to get started.  While most of the strongly-typed resources have complete examples, the examples for loosely-typed resources like `powerplatform_data_record` only cover a fraction of the possible configuration options.  See the `/examples/resources/powerplatform_data_record` folder for examples and feel free to suggest and contribute more.

Examples of real-world use cases are encouraged.  Please contribute those types of examples to the [Power Platform Terraform QuickStarts](https://github.com/microsoft/power-platform-terraform-quickstarts) repo.

### Data Sources

Creating a new [data source](https://developer.hashicorp.com/terraform/plugin/framework/data-sources) can allow terraform to read reference data about Power Platform services and infrastrucutre.  Implementing a data source will help you learn some of the concepts that will be useful in eventually developing a resource, but data sources are much simpler since you only have to handle reading data.  The issue backlog contains a [list of proposed data sources](https://github.com/microsoft/terraform-provider-power-platform/issues?q=is%3Aissue%20state%3Aopen%20label%3A%22data%20source%22) that may need a contributor.

### Resources

Creating a new [resource](https://developer.hashicorp.com/terraform/plugin/framework/resources) can allow terraform to manage new Power Platform infrastructure not currently provided by the provider.  Resources are the most complex since you need to implement the full resource lifecycle.  Feel free to add a comment to the issue if you'd like to start a conversation on making a contribution for that resource request. The issue backlog contains a [list of proposed resources](https://github.com/microsoft/terraform-provider-power-platform/issues?q=is%3Aissue%20state%3Aopen%20label%3Aresource) that may need a contributor.

## Pull Request Checklist

PRs for new resources or data sources are expected to meeting the following criteria:

- Add a production quality implementation of the resource or datasource in [/internal](/internal/)
- Add unit tests and acceptance tests for your contribution in [/internal](/internal/).  
  - Tests should pass and provide >80% coverage of your contribution
- Add examples for your contribution in [/examples](/examples/) (see [Terraform Documentaion on examples](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation#add-configuration-examples))
- Add [schema descriptions](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation#add-schema-descriptions) for your resource or data source in [/internal](/internal/)
- and/or [/templates](/templates/)
- [Update auto-generated documentation](./DEVELOPER.md#updating-documentation) in [/docs](/docs/). (Don't manually edit [/docs](/docs/) or your updates will be overwritten)
- Ensure the PR description clearly describes the feature you're adding and any known limitations (We recommend using GitHub Copilot for PR descriptions)

## Getting Started

Ready to contribute? Check out our [Developer Guide](./DEVELOPER.md) for detailed instructions on setting up your development environment, building the provider, and running tests. We look forward to your contributions!
