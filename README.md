# Terraform Provider for Phillips Hue

My Phillips Hue Hub is configured using this provider, it should be solid enough for most applications.

This is my first terraform provider and I did it while learning "go".
I also did a lot of experimentation on the Hue Bridge REST API but that was the easiest part.

Any suggestion/contribution is appreciated, especially if it is regarding tests and best practices in terraform and go.

## How to develop

This is the way I did:

- I am on MacOS, if you are using another platform you will need to adjust the `OS_ARCH` in the `Makefile`

- Created a folder called `definitions` (actually I `ln -s` another path), you can use `examples` as a base.

- To generate username (see `examples/provider.tf`), press the link button in the bridge then:

```sh
curl -X POST http://192.168.xxx.yyy/api -d '{"devicetype":"hue#terraform"}'
```

- Created a `~/.terraformrc` to use the plugin under development, ex (`/xxx/yyy/go` is the `GOPATH`):

```
provider_installation {
  dev_overrides {
      "github.com/roperto/hue" = "/xxx/yyy/go/bin"
      }
  direct {}
}
```

- Use `make apply` -- notice it has flags to enable debugging.

## Publish version

I did not bother publishing it to terraform repository as I don't know if anyone else
will ever use it -- if I see interest I will try to get it published properly.

In my local usage, I just apply with the development options.

Feel free to open an issue if you want more details.
