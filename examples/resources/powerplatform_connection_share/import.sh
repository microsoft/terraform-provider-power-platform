# A connection share is the permission of one principal on one connection, so it is imported using the composite id
# {environment_id}/{connector_name}/{connection_id}/{principal_entra_object_id}
# 00000000-0000-0000-0000-000000000001 = environment id
# shared_sql                           = connector api name (the `name` of the powerplatform_connection resource)
# 00000000000000000000000000000002     = connection id (the `id` of the powerplatform_connection resource)
# 00000000-0000-0000-0000-000000000003 = Entra object id of the principal the connection is shared with
terraform import powerplatform_connection_share.example "00000000-0000-0000-0000-000000000001/shared_sql/00000000000000000000000000000002/00000000-0000-0000-0000-000000000003"
