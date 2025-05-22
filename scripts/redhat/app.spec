Name: cloud-manager-server
Version: %VERSION%
Release: 0
License: MIT
Requires: openssl
Prefix: /opt
Summary: example go-fiber application template for aciklab apps
Group: Applications/System
BuildArch: x86_64

%description
example go-fiber application template for aciklab apps

%pre

%prep

%build

%install
mkdir -p %{buildroot}/opt/cloud-manager-server
cp -rfa %{_sourcedir}/cloud-manager-server/* %{buildroot}/opt/cloud-manager-server/
cp -a %{_sourcedir}/cloud-manager-server/.env.example %{buildroot}/opt/cloud-manager-server/


%post -p /bin/bash
if ! getent passwd cloud-manager-server >/dev/null; then
    adduser --system --user-group \
        --home-dir /run/cloud-manager-server \
        --shell /bin/bash \
        cloud-manager-server
fi

if ! [ -d "/run/cloud-manager-server" ]; then
    mkdir -p "/run/cloud-manager-server"
    chown cloud-manager-server:cloud-manager-server "/run/cloud-manager-server"
fi
if ! [ -f "/opt/cloud-manager-server/keys/cloud-manager-server.key" ]; then
    mkdir -p '/opt/cloud-manager-server/keys/'
    openssl req -x509 -newkey rsa:4096 -subj "/CN=$(hostname -I | cut -d" " -f1 | xargs)" -extensions SAN -reqexts SAN -config <(cat $(echo "$(openssl version -d | sed 's/.*"\(.*\)"/\1/g')/openssl.cnf") <(printf "\n[SAN]\nsubjectAltName=IP:$(hostname -I | cut -d" " -f1 | xargs),IP:127.0.0.1,DNS:$(hostname)")) -keyout /opt/cloud-manager-server/keys/cloud-manager-server.key -nodes -out /opt/cloud-manager-server/keys/cloud-manager-server.pem -sha256 -days 358000
fi

chown -R cloud-manager-server:cloud-manager-server /opt/cloud-manager-server
chmod -R 770 /opt/cloud-manager-server
if [ -f "/usr/lib/systemd/system/cloud-manager-server.service" ]; then
    rm -rf /usr/lib/systemd/system/cloud-manager-server.service
    systemctl disable cloud-manager-server.service
    systemctl stop cloud-manager-server.service
    systemctl daemon-reload cloud-manager-server.service
fi

echo """
[Unit]
Description=Fiber App Template %I

[Service]
Type=simple
WorkingDirectory=/opt/cloud-manager-server
ExecStart=/opt/cloud-manager-server/app -type=%i
Restart=always
RestartSec=10
KillSignal=SIGINT
SyslogIdentifier=cloud-manager-server
User=root
Group=root

[Install]
WantedBy=multi-user.target
    """ > /etc/systemd/system/cloud-manager-server@.service

systemctl daemon-reload

%clean

%files
%defattr(0770, root, root)
/opt/cloud-manager-server/app
/opt/cloud-manager-server/.env.example

%define _unpackaged_files_terminate_build 0