---
subcategory: "Elastic Load Balance (ELB)"
---

# flexibleengine_lb_pool_v2

Manages an **enhanced** load balancer pool resource within FlexibleEngine.

## Example Usage

```hcl
resource "flexibleengine_lb_pool_v2" "pool_1" {
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  persistence {
    type        = "HTTP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the V2 Networking client.
    A Networking client is needed to create an . If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    pool.

* `name` - (Optional) Human-readable name for the pool.

* `description` - (Optional) Human-readable description for the pool.

* `protocol` = (Required) The protocol - can either be TCP, UDP or HTTP.

    - When the protocol used by the listener is UDP, the protocol of the backend pool must be UDP.
    - When the protocol used by the listener is TCP, the protocol of the backend pool must be TCP.
    - When the protocol used by the listener is HTTP or TERMINATED_HTTPS, the protocol of the backend pool must be HTTP.

    Changing this creates a new pool.

* `loadbalancer_id` - (Optional) The load balancer on which to provision this
    pool. Changing this creates a new pool.
    Note:  One of LoadbalancerID or ListenerID must be provided.

* `listener_id` - (Optional) The Listener on which the members of the pool
    will be associated with. Changing this creates a new pool.
	Note:  One of LoadbalancerID or ListenerID must be provided.

* `lb_method` - (Required) The load balancing algorithm to
    distribute traffic to the pool's members. Must be one of
    ROUND_ROBIN, LEAST_CONNECTIONS, or SOURCE_IP.

* `persistence` - Omit this field to prevent session persistence.  Indicates
    whether connections in the same session will be processed by the same Pool
    member or not. Changing this creates a new pool.

* `tenant_id` - (Optional) The UUID of the tenant who owns the pool.
    Only administrative users can specify a tenant UUID other than their own.
    Changing this creates a new pool.

* `admin_state_up` - (Optional) The administrative state of the pool.
    A valid value is true (UP) or false (DOWN).

The `persistence` argument supports:

* `type` - (Required) The type of persistence mode. The current specification
    supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE.

* `cookie_name` - (Optional) The name of the cookie if persistence mode is set
    appropriately. Required if `type = APP_COOKIE`.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the pool.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `protocol` - See Argument Reference above.
* `lb_method` - See Argument Reference above.
* `persistence` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
