make build_linux 
sudo systemctl stop goi_example.service
sudo cp build/goi_example /opt/goi_example/
sudo systemctl restart goi_example.service