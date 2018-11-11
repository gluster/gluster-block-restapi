%if 0%{?fedora}
%global with_bundled 1
%else
%global with_bundled 1
%endif

%{!?with_debug: %global with_debug 1}

%if 0%{?with_debug}
%global _dwz_low_mem_die_limit 0
%else
%global debug_package   %{nil}
%endif

%{!?go_arches: %global go_arches x86_64 aarch64 ppc64le }

%global provider github
%global provider_tld com
%global project aravindavk
%global repo gluster-block-restapi
%global provider_prefix %{provider}.%{provider_tld}/%{project}/%{repo}
%global import_path %{provider_prefix}

%global gluster_block_restapi_make %{__make} PREFIX=%{_prefix} EXEC_PREFIX=%{_exec_prefix} BINDIR=%{_bindir} SBINDIR=%{_sbindir} DATADIR=%{_datadir} LOCALSTATEDIR=%{_sharedstatedir} LOGDIR=%{_localstatedir}/log SYSCONFDIR=%{_sysconfdir} FASTBUILD=off

%global gluster_block_restapi_ver 1
%global gluster_block_restapi_rel 0

Name: %{repo}
Version: %{gluster_block_restapi_ver}
Release: %{gluster_block_restapi_rel}%{?dist}
Summary: REST APIs for Gluster Block Volume management
License: GPLv2 or LGPLv3+
URL: https://%{provider_prefix}
%if 0%{?with_bundled}
Source0: https://%{provider_prefix}/releases/download/v%{version}/gluster-block-restapi-v%{gluster_block_restapi_ver}-%{gluster_block_restapi_rel}-vendor.tar.xz
%else
Source0: https://%{provider_prefix}/releases/download/v%{version}/gluster-block-restapi-v%{gluster_block_restapi_ver}-%{gluster_block_restapi_rel}.tar.xz
%endif

ExclusiveArch: %{go_arches}

BuildRequires: %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}
BuildRequires: systemd

%if ! 0%{?with_bundled}
BuildRequires: golang(github.com/BurntSushi/toml)
BuildRequires: golang(github.com/sirupsen/logrus)
BuildRequires: golang(github.com/gorilla/mux)
%endif

Requires: /usr/bin/strings
%{?systemd_requires}

%description
The project gluster-block-restapi provides set of APIs for Gluster
block volumes management.

%prep
%setup -q -n gluster-block-restapi

%build
export GOPATH=$(pwd):%{gopath}
mkdir -p src/%(dirname %{import_path})
ln -s ../../../ src/%{import_path}

pushd src/%{import_path}
# Build gluster-block-restapi
%{gluster_block_restapi_make} glusterblockrestd
popd

%install
# Install gluster-block-restapi
%{gluster_block_restapi_make} DESTDIR=%{buildroot} install

%post
%systemd_post glusterblockrestd.service

%preun
%systemd_preun glusterblockrestd.service

%files
%{_sbindir}/glusterblockrestd
%{_unitdir}/glusterblockrestd.service
%{_sysconfdir}/gluster-block-restapi/config.toml

%changelog
* Sun Nov 11 2018 Aravinda VK <avishwan@redhat.com> - 1.0.0-1
- Initial Spec
