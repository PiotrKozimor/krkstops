# sudo dnf install coreos-installer
# stream="stable"
# coreos-installer download --decompress -s "${stream}" -p azure -f vhd.xz
set -xe
az_region="germanywestcentral"
az_resource_group="coreos"
az_storage_account="pcpiotr1"
az_container="coreos"
downloaded_image_file="~/Downloads/fedora-coreos-34.20210808.3.0-azure.x86_64.vhd"
az_image_name="my-fcos-image"
az_image_blob="${az_image_name}.vhd"
# az group create -l "${az_region}" -n "${az_resource_group}"
# az storage account create -g "${az_resource_group}" -n "${az_storage_account}"
# cs=$(az storage account show-connection-string -n "${az_storage_account}" | jq -r .connectionString)
# az storage container create --connection-string "${cs}" -n "${az_container}"
# az storage blob upload --connection-string "${cs}" -c "${az_container}" -f "${downloaded_image_file}" -n "${az_image_blob}"
# az image create -n "${az_image_name}" -g "${az_resource_group}" --source "https://${az_storage_account}.blob.core.windows.net/${az_container}/${az_image_blob}" --location "${az_region}" --os-type Linux
# az storage blob delete --connection-string "$cs" -c "${az_container}" -n "${az_image_blob}"
ignition_path="./butane.ign"
az vm create \
    -n krkstops \
    -g "${az_resource_group}" \
    --image "${az_image_name}" \
    --admin-username core \
    --size Standard_B1s \
    --storage-sku StandardSSD_LRS \
    --os-disk-size-gb 64 \
    --public-ip-sku Standard \
    --custom-data "$(cat ${ignition_path})"
az vm open-port --resource-group "${az_resource_group}" --name krkstops --port 8080
az vm open-port --resource-group "${az_resource_group}" --name krkstops --priority 800 --port 9090

