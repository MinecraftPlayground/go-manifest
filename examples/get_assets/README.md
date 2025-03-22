# Download asset file

This example downloads a specified asset file for a specified version. And saves it as file in `assets/<asset_id>`

### Usage

`main.go <version_id> <asset_id>` where `version_id` is the version you want it from and `asset_id` the file path of the asset file you want to download.

### Example
```sh
go run main.go 1.21.4 minecraft/sounds/damage/fallbig.ogg
```
