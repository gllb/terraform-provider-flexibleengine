---
subcategory: "Elastic Cloud Server (ECS)"
---

# flexibleengine\_compute\_instance_v2

Manages a V2 VM instance resource within FlexibleEngine.

## Example Usage

### Basic Instance

```hcl
resource "flexibleengine_compute_instance_v2" "basic" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }

  tags = {
    foo  = "bar"
    this = "that"
  }
}
```

### Instance With Attached Volume

```hcl
resource "flexibleengine_blockstorage_volume_v2" "myvol" {
  name = "myvol"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "myinstance" {
  name            = "myinstance"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
}

resource "flexibleengine_compute_volume_attach_v2" "attached" {
  instance_id = flexibleengine_compute_instance_v2.myinstance.id
  volume_id   = flexibleengine_blockstorage_volume_v2.myvol.id
}
```

### Boot From Volume

```hcl
resource "flexibleengine_compute_instance_v2" "boot-from-volume" {
  name            = "boot-from-volume"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    volume_size           = 5
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
    volume_type           = "SSD"
  }

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
}
```

### Boot From an Existing Volume

```hcl
resource "flexibleengine_blockstorage_volume_v2" "myvol" {
  name     = "myvol"
  size     = 5
  image_id = "<image-id>"
}

resource "flexibleengine_compute_instance_v2" "boot-from-volume" {
  name            = "bootfromvolume"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = flexibleengine_blockstorage_volume_v2.myvol.id
    source_type           = "volume"
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
}
```

### Boot Instance, Create Volume, and Attach Volume as a Block Device

```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    source_type           = "blank"
    destination_type      = "volume"
    volume_size           = 1
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Boot Instance and Attach Existing Volume as a Block Device

```hcl
resource "flexibleengine_blockstorage_volume_v2" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "instance_1"
  image_id        = "<image-id>"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    uuid                  = "<image-id>"
    source_type           = "image"
    destination_type      = "local"
    boot_index            = 0
    delete_on_termination = true
  }

  block_device {
    uuid                  = flexibleengine_blockstorage_volume_v2.volume_1.id
    source_type           = "volume"
    destination_type      = "volume"
    boot_index            = 1
    delete_on_termination = true
  }
}
```

### Instance With Multiple Networks

```hcl
resource "flexibleengine_networking_floatingip_v2" "myip" {
  pool = "admin_external_net"
}

resource "flexibleengine_compute_instance_v2" "multi-net" {
  name            = "multi-net"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }

  network {
    uuid = "55534eaa-533a-419d-9b40-ec427ea7195a"
  }
}

resource "flexibleengine_compute_floatingip_associate_v2" "myip" {
  floating_ip = flexibleengine_networking_floatingip_v2.myip.address
  instance_id = flexibleengine_compute_instance_v2.multi-net.id
  fixed_ip    = flexibleengine_compute_instance_v2.multi-net.network.1.fixed_ip_v4
}
```

### Instance with Multiple Ephemeral Disks

```hcl
resource "flexibleengine_compute_instance_v2" "multi-eph" {
  name            = "multi_eph"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image-id>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instance with User Data (cloud-init)

```hcl
resource "flexibleengine_compute_instance_v2" "instance_1" {
  name            = "basic"
  image_id        = "ad091b52-742f-469e-8f3c-fd81cadf0743"
  flavor_id       = "s3.large.2"
  key_pair        = "my_key_pair_name"
  security_groups = ["default"]
  user_data       = "#cloud-config\nhostname: instance_1.example.com\nfqdn: instance_1.example.com"

  network {
    uuid = "3c4a0d74-24b9-46cf-9d7f-8b7a4dc2f65c"
  }
}
```

`user_data` can come from a variety of sources: inline, read in from the `file`
function, or the `template_cloudinit_config` resource.

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the server instance. If
    omitted, the `region` argument of the provider is used. Changing this
    creates a new server.

* `name` - (Required) A unique name for the resource.

* `image_id` - (Optional; Required if `image_name` is empty and not booting
    from a volume. Do not specify if booting from a volume.) The image ID of
    the desired image for the server. Changing this creates a new server.

* `image_name` - (Optional; Required if `image_id` is empty and not booting
    from a volume. Do not specify if booting from a volume.) The name of the
    desired image for the server. Changing this creates a new server.

* `flavor_id` - (Optional; Required if `flavor_name` is empty) The flavor ID of
    the desired flavor for the server. Changing this resizes the existing server.

* `flavor_name` - (Optional; Required if `flavor_id` is empty) The name of the
    desired flavor for the server. Changing this resizes the existing server.

* `availability_zone` - (Optional) The availability zone in which to create
    the server. Changing this creates a new server.

* `security_groups` - (Optional) An array of one or more security group names
    to associate with the server. Changing this results in adding/removing
    security groups from the existing server. *Note*: When attaching the
    instance to networks using Ports, place the security groups on the Port
    and not the instance.

* `network` - (Optional) An array of one or more networks to attach to the
    instance. The network object structure is documented below. Changing this
    creates a new server.

* `user_data` - (Optional) The user data to provide when launching the instance.
    Changing this creates a new server.

* `metadata` - (Optional) The key/value pairs to associate with the instance.
    Changing this updates the existing server metadata.

* `config_drive` - (Optional) Whether to use the config_drive feature to
    configure the instance. Changing this creates a new server.

* `admin_pass` - (Optional) The administrative password to assign to the server.
    Changing this changes the root password on the existing server.

* `key_pair` - (Optional) The name of a key pair to put on the server. The key
    pair must already be created and associated with the tenant's account.
    Changing this creates a new server.

* `block_device` - (Optional) Configuration of block devices. The block_device
    structure is documented below. Changing this creates a new server.
    You can specify multiple block devices which will create an instance with
    multiple disks. This configuration is very flexible, so please see the
    following [reference](http://docs.openstack.org/developer/nova/block_device_mapping.html)
    for more information.

* `scheduler_hints` - (Optional) Provide the Nova scheduler with hints on how
    the instance should be launched. The available hints are described below.

* `stop_before_destroy` - (Optional) Whether to try stop instance gracefully
    before destroying it, thus giving chance for guest OS daemons to stop correctly.
    If instance doesn't stop within timeout, it will be destroyed anyway.

* `force_delete` - (Optional) Whether to force the FlexibleEngine instance to be
    forcefully deleted. This is useful for environments that have reclaim / soft
    deletion enabled.

* `auto_recovery` - (Optional) Configures or deletes automatic recovery of an instance

* `tags` - (Optional) Tags key/value pairs to associate with the instance.

The `network` block supports:

* `uuid` - (Required unless `port` is provided) The network UUID to
    attach to the server. Changing this creates a new server.

* `port` - (Required unless `uuid` is provided) The port UUID of a network to
    attach to the server. Changing this creates a new server.

* `fixed_ip_v4` - (Optional) Specifies a fixed IPv4 address to be used on this
    network. Changing this creates a new server.

* `fixed_ip_v6` - (Optional) Specifies a fixed IPv6 address to be used on this
    network. Changing this creates a new server.

* `access_network` - (Optional) Specifies if this network should be used for
    provisioning access. Accepts true or false. Defaults to false.

The `block_device` block supports:

* `uuid` - (Required unless `source_type` is set to `"blank"` ) The UUID of
    the image, volume, or snapshot. Changing this creates a new server.

* `source_type` - (Required) The source type of the device. Must be one of
    "blank", "image", "volume", or "snapshot". Changing this creates a new
    server.

* `volume_size` - The size of the volume to create (in gigabytes). Required
    in the following combinations: source=image and destination=volume,
    source=blank and destination=local, and source=blank and destination=volume.
    Changing this creates a new server.

* `volume_type` - (Optional) Currently, the value can be `SSD` (ultra-I/O disk type),
    `SAS` (high I/O disk type), or `SATA` (common I/O disk type)

* `boot_index` - (Optional) The boot index of the volume. It defaults to 0, which
	indicates that it's a system disk. Changing this creates a new server.

* `destination_type` - (Optional) The type that gets created. Possible values
    are "volume" and "local". Changing this creates a new server.

* `delete_on_termination` - (Optional) Delete the volume / block device upon
    termination of the instance. Defaults to false. Changing this creates a
    new server.

* `disk_bus` - (Optional) The low-level disk bus that will be used, for example, *virtio*, *scsi*.
    Most common thing is to leave this empty. Changing this creates a new server.

The `scheduler_hints` block supports:

* `group` - (Optional) A UUID of a Server Group. The instance will be placed
    into that group.

* `different_host` - (Optional) A list of instance UUIDs. The instance will
    be scheduled on a different host than all other instances.

* `same_host` - (Optional) A list of instance UUIDs. The instance will be
    scheduled on the same host of those specified.

* `query` - (Optional) A conditional query that a compute node must pass in
    order to host an instance.

* `target_cell` - (Optional) The name of a cell to host the instance.

* `build_near_host_ip` - (Optional) An IP Address in CIDR form. The instance
    will be placed on a compute node that is in the same subnet.

## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `access_ip_v4` - The first detected Fixed IPv4 address _or_ the Floating IP.
* `access_ip_v6` - The first detected Fixed IPv6 address.
* `metadata` - See Argument Reference above.
* `security_groups` - See Argument Reference above.
* `flavor_id` - See Argument Reference above.
* `flavor_name` - See Argument Reference above.
* `network/uuid` - See Argument Reference above.
* `network/port` - See Argument Reference above.
* `network/name` - The human-readable name of the network.
* `network/fixed_ip_v4` - The Fixed IPv4 address of the Instance on that network.
* `network/fixed_ip_v6` - The Fixed IPv6 address of the Instance on that network.
* `network/mac` - The MAC address of the NIC on that network.
* `all_metadata` - Contains all instance metadata, even metadata not set
    by Terraform.
* `auto_recovery` - See Argument Reference above.
* `tags` - See Argument Reference above.
* `floating_ip` - The EIP address that is associted to the instance.
* `system_disk_id` - The system disk voume ID.
* `volume_attached` - An array of one or more disks to attach to the instance.
    The volume_attached object structure is documented below.
* `status` - The status of the instance.

The `volume_attached` block supports:

* `uuid` - The volume id on that attachment.
* `boot_index` - The volume boot index on that attachment.
* `size` - The volume size on that attachment.
* `type` - The volume type on that attachment.
* `pci_address` - The volume pci address on that attachment.

## Import

Instances can be imported by their `id`. For example,
```
terraform import flexibleengine_compute_instance_v2.my_instance b11b407c-e604-4e8d-8bc4-92398320b847
```
Note that the imported state may not be identical to your resource definition, due to some attrubutes
missing from the API response, security or some other reason. The missing attributes include:
`admin_pass`, `config_drive`, `user_data`, `block_device`, `scheduler_hints`, `stop_before_destroy`,
`network/access_network` and arguments for pre-paid. It is generally recommended running
`terraform plan` after importing an instance. You can then decide if changes should
be applied to the instance, or the resource definition should be updated to align
with the instance. Also you can ignore changes as below.
```
resource "flexibleengine_compute_instance_v2" "my_instance" {
    ...

  lifecycle {
    ignore_changes = [
      user_data, block_device,
    ]
  }
}
```

## Notes

### Multiple Ephemeral Disks

It's possible to specify multiple `block_device` entries to create an instance
with multiple ephemeral (local) disks. In order to create multiple ephemeral
disks, the sum of the total amount of ephemeral space must be less than or
equal to what the chosen flavor supports.

The following example shows how to create an instance with multiple ephemeral
disks:

```
resource "flexibleengine_compute_instance_v2" "foo" {
  name            = "terraform-test"
  security_groups = ["default"]

  block_device {
    boot_index            = 0
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "image"
    uuid                  = "<image uuid>"
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }

  block_device {
    boot_index            = -1
    delete_on_termination = true
    destination_type      = "local"
    source_type           = "blank"
    volume_size           = 1
  }
}
```

### Instances and Ports

Neutron Ports are a great feature and provide a lot of functionality. However,
there are some notes to be aware of when mixing Instances and Ports:

* When attaching an Instance to one or more networks using Ports, place the
security groups on the Port and not the Instance. If you place the security
groups on the Instance, the security groups will not be applied upon creation,
but they will be applied upon a refresh. This is a known FlexibleEngine bug.

* Network IP information is not available within an instance for networks that
are attached with Ports. This is mostly due to the flexibility Neutron Ports
provide when it comes to IP addresses. For example, a Neutron Port can have
multiple Fixed IP addresses associated with it. It's not possible to know which
single IP address the user would want returned to the Instance's state
information. Therefore, in order for a Provisioner to connect to an Instance
via it's network Port, customize the `connection` information:

```hcl
resource "flexibleengine_networking_port_v2" "port_1" {
  name           = "port_1"
  admin_state_up = "true"
  network_id     = "0a1d0a27-cffa-4de3-92c5-9d3fd3f2e74d"

  security_group_ids = [
    "2f02d20a-8dca-49b7-b26f-b6ce9fddaf4f",
    "ca1e5ed7-dae8-4605-987b-fadaeeb30461",
  ]
}

resource "flexibleengine_compute_instance_v2" "instance_1" {
  name = "instance_1"

  network {
    port = flexibleengine_networking_port_v2.port_1.id
  }

  connection {
    user        = "root"
    host        = flexibleengine_networking_port_v2.port_1.fixed_ip.0.ip_address
    private_key = "~/path/to/key"
  }

  provisioner "remote-exec" {
    inline = [
      "echo terraform executed > /tmp/foo",
    ]
  }
}
```
