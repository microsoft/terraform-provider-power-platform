# Connection resource can be imported using the composite id {environment_id}/{connector_name}/{connection_id}
# (replace with a real environment id, the connector api name (for example shared_azureopenai) and the connection id
# that can be found in the URL of the connection details page in the Power Apps maker portal)
terraform import powerplatform_connection.example 00000000-0000-0000-0000-000000000000/shared_azureopenai/00000000000000000000000000000000
