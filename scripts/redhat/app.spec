Name: fiber-app-template
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
mkdir -p %{buildroot}/opt/fiber-app-template
cp -rfa %{_sourcedir}/fiber-app-template/* %{buildroot}/opt/fiber-app-template/
cp -a %{_sourcedir}/fiber-app-template/.env.example %{buildroot}/opt/fiber-app-template/


%post -p /bin/bash
if ! getent passwd fiber-app-template >/dev/null; then
    adduser --system --user-group \
        --home-dir /run/fiber-app-template \
        --shell /bin/bash \
        fiber-app-template
fi

if ! [ -d "/run/fiber-app-template" ]; then
    mkdir -p "/run/fiber-app-template"
    chown fiber-app-template:fiber-app-template "/run/fiber-app-template"
fi
if ! [ -f "/opt/fiber-app-template/keys/fiber-app-template.key" ]; then
    mkdir -p '/opt/fiber-app-template/keys/'
    openssl req -x509 -newkey rsa:4096 -subj "/CN=$(hostname -I | cut -d" " -f1 | xargs)" -extensions SAN -reqexts SAN -config <(cat $(echo "$(openssl version -d | sed 's/.*"\(.*\)"/\1/g')/openssl.cnf") <(printf "\n[SAN]\nsubjectAltName=IP:$(hostname -I | cut -d" " -f1 | xargs),IP:127.0.0.1,DNS:$(hostname)")) -keyout /opt/fiber-app-template/keys/fiber-app-template.key -nodes -out /opt/fiber-app-template/keys/fiber-app-template.pem -sha256 -days 358000
fi

chown -R fiber-app-template:fiber-app-template /opt/fiber-app-template
chmod -R 770 /opt/fiber-app-template
if [ -f "/usr/lib/systemd/system/fiber-app-template.service" ]; then
    rm -rf /usr/lib/systemd/system/fiber-app-template.service
    systemctl disable fiber-app-template.service
    systemctl stop fiber-app-template.service
    systemctl daemon-reload fiber-app-template.service
fi

echo """
[Unit]
Description=Fiber App Template %I

[Service]
Type=simple
WorkingDirectory=/opt/fiber-app-template
ExecStart=/opt/fiber-app-template/app -type=%i
Restart=always
RestartSec=10
KillSignal=SIGINT
SyslogIdentifier=fiber-app-template
User=root
Group=root

[Install]
WantedBy=multi-user.target
    """ > /etc/systemd/system/fiber-app-template@.service

systemctl daemon-reload

%clean

%files
%defattr(0770, root, root)
/opt/fiber-app-template/app
/opt/fiber-app-template/.env.example

%define _unpackaged_files_terminate_build 0