Name:           overthink
# Version is injected at build time via: rpmbuild --define "package_version X.Y.Z"
%{!?package_version: %define package_version 0.0.0}
Version:        %{package_version}
Release:        1%{?dist}
Summary:        A dramatic overanalysis engine for your questionable life decisions

License:        MIT
URL:            https://github.com/rishichawda/overthinker
Source0:        https://github.com/rishichawda/overthinker/archive/refs/tags/v%{package_version}.tar.gz

BuildRequires:  golang >= 1.24

%description
overthink answers your questions in the most unnecessarily elaborate way
possible — complete with fabricated statistics, fictional citations,
an Emotional Risk Index, and ASCII bar charts.

%prep
%autosetup -n overthinker-%{version}

%build
export CGO_ENABLED=0
go build -trimpath -ldflags "-s -w" -o overthink ./cmd/overthink

%install
install -Dpm 0755 overthink %{buildroot}%{_bindir}/overthink

%files
%license LICENSE
%{_bindir}/overthink

%changelog
* Thu Feb 26 2026 Rishi Chawda <rishichawda@users.noreply.github.com> - 1.0.0-1
- Initial release
