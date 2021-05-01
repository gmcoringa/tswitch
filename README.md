# Terragrunt and Terraform switcher [![push](https://github.com/gmcoringa/tswitch/workflows/push/badge.svg?branch=main&event=push)](https://github.com//gmcoringa/tswitch/actions)

The `tswitch` command line tool switches between different versions of [terraform](https://www.terraform.io/) and [terragrunt](https://terragrunt.gruntwork.io/) looking for constraints in `terragrunt.hcl` file. If you wish to use other methods, there are other switches that you may use. 

## Installation

1. Go to the [Releases Page](https://github.com/gmcoringa/tswitch/releases).
2. Downloading the binary for your operating system: e.g., if you’re on a Mac, download tswitch_darwin_amd64; if you’re on Windows, download tswitch_windows_amd64.exe, etc.
3. Rename the downloaded file to `tswitch`.
4. Add execute permissions to the binary. E.g., On Linux and Mac: `chmod u+x tswitch`.
5. Put the binary somewhere on your PATH. E.g., On Linux and Mac: `mv tswitch /usr/local/bin/tswitch`.

## Usage

Just point to a terragrunt file, ex:
```bash
tswitch -terragrunt /some/path/terragrunt.hcl
# Help to see all available options
tswitch -help
```

Your `terragrunt.hcl` **must** have the constraints declared, ex:
```
terragrunt_version_constraint = ">= 0.26, < 0.27"
terraform_version_constraint  = ">= 0.13, < 0.14"
```

## Configuration file

All command line options are also supported by a configuration file, ex:
```yaml
terragruntFile: ./terragrunt.hcl
installDir: /tmp/tswitch/bin
cacheDir: /tmp/tsiwtch/data
```

Use the flag `-config` to set the path for a configuration file.

## Issues

Please open  *issues* here:  [New Issue](https://github.com/gmcoringa/tswitch/issues)

## License

This code is released under the GNU GENERAL PUBLIC LICENSE. See [LICENSE](LICENSE).
