Name: cloud-manager
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
mkdir -p %{buildroot}/opt/cloud-manager
cp -rfa %{_sourcedir}/cloud-manager/* %{buildroot}/opt/cloud-manager/
cp -a %{_sourcedir}/cloud-manager/.env.example %{buildroot}/opt/cloud-manager/


%post -p /bin/bash
if ! getent passwd cloud-manager >/dev/null; then
    adduser --system --user-group \
        --home-dir /run/cloud-manager \
        --shell /bin/bash \
        cloud-manager
fi

if ! [ -d "/run/cloud-manager" ]; then
    mkdir -p "/run/cloud-manager"
    chown cloud-manager:cloud-manager "/run/cloud-manager"
fi
if ! [ -f "/opt/cloud-manager/keys/cloud-manager.key" ]; then
    mkdir -p '/opt/cloud-manager/keys/'
    openssl req -x509 -newkey rsa:4096 -subj "/CN=$(hostname -I | cut -d" " -f1 | xargs)" -extensions SAN -reqexts SAN -config <(cat $(echo "$(openssl version -d | sed 's/.*"\(.*\)"/\1/g')/openssl.cnf") <(printf "\n[SAN]\nsubjectAltName=IP:$(hostname -I | cut -d" " -f1 | xargs),IP:127.0.0.1,DNS:$(hostname)")) -keyout /opt/cloud-manager/keys/cloud-manager.key -nodes -out /opt/cloud-manager/keys/cloud-manager.pem -sha256 -days 358000
fi

chown -R cloud-manager:cloud-manager /opt/cloud-manager
chmod -R 770 /opt/cloud-manager
if [ -f "/usr/lib/systemd/system/cloud-manager.service" ]; then
    rm -rf /usr/lib/systemd/system/cloud-manager.service
    systemctl disable cloud-manager.service
    systemctl stop cloud-manager.service
    systemctl daemon-reload cloud-manager.service
fi

echo """
[Unit]
Description=Fiber App Template %I

[Service]
Type=simple
WorkingDirectory=/opt/cloud-manager
ExecStart=/opt/cloud-manager/app -type=%i
Restart=always
RestartSec=10
KillSignal=SIGINT
SyslogIdentifier=cloud-manager
User=root
Group=root

[Install]
WantedBy=multi-user.target
    """ > /etc/systemd/system/cloud-manager@.service

systemctl daemon-reload

%clean

%files
%defattr(0770, root, root)
/opt/cloud-manager/app
/opt/cloud-manager/.env.example

%define _unpackaged_files_terminate_build 0