#!/bin/sh

systemctl --user stop uapd.service 2>/dev/null

sudo cp -v ./uapd /usr/local/bin/

install -v ./uapd.service ~/.config/systemd/user/uapd.service

systemctl --user enable --now uapd.service

systemctl --user status uapd


