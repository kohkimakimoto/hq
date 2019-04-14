Name:           %{_product_name}
Version:        %{_product_version}

Release:        1.el%{_rhel_version}
Summary:        Simplistic job queue engine
Group:          Development/Tools
License:        MIT
Source0:        %{name}_linux_amd64.zip
Source1:        hq.toml
Source2:        hq.sysconfig
Source3:        hq.service
BuildRoot:      %(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description
Server Management Framework

%prep
%setup -q -c

%install
mkdir -p %{buildroot}/%{_bindir}
cp %{name} %{buildroot}/%{_bindir}

mkdir -p %{buildroot}/%{_sysconfdir}/%{name}
cp %{SOURCE1} %{buildroot}/%{_sysconfdir}/hq/hq.toml

mkdir -p %{buildroot}/%{_sysconfdir}/sysconfig
cp %{SOURCE2} %{buildroot}/%{_sysconfdir}/sysconfig/hq

mkdir -p %{buildroot}/var/lib/hq

%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
mkdir -p %{buildroot}/%{_unitdir}
cp %{SOURCE3} %{buildroot}/%{_unitdir}/
%endif

%pre
getent group hq >/dev/null || groupadd -r hq
getent passwd hq >/dev/null || \
    useradd -r -g hq -d /var/lib/hq -s /sbin/nologin \
    -c "hq user" hq
exit 0

%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
%post
%systemd_post hq.service
systemctl daemon-reload

%preun
%systemd_preun hq.service
systemctl daemon-reload

%endif

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%attr(755, root, root) %{_bindir}/%{name}
%dir %attr(755, hq, hq) /var/lib/hq
%config(noreplace) %{_sysconfdir}/hq/hq.toml
%config(noreplace) %{_sysconfdir}/sysconfig/hq

%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
%{_unitdir}/hq.service
%endif

%doc
