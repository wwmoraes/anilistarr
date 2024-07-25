let
  ## nix-prefetch-url --unpack <url>
  nixpkgs = fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/refs/tags/24.05.tar.gz";
    sha256 = "1lr1h35prqkd1mkmzriwlpvxcb34kmhc9dnr48gkm8hh089hifmx";
  };
  nixpkgs-unstable = fetchTarball {
    name = "nixos-unstable-a14c5d651cee9ed70f9cd9e83f323f1e531002db";
    # url = "https://github.com/NixOS/nixpkgs/archive/refs/heads/nixpkgs-unstable.tar.gz";
    url = "https://github.com/NixOS/nixpkgs/archive/a14c5d651cee9ed70f9cd9e83f323f1e531002db.tar.gz";
    sha256 = "1b2dwbqm5vdr7rmxbj5ngrxm7sj5r725rqy60vnlirbbwks6aahb";
  };
  nur = fetchTarball {
    url = "https://github.com/nix-community/NUR/archive/master.tar.gz";
    sha256 = "08c7qaxw2mg4rj80w571smnrf152km1axw639dsl5idhhc9bgqhm";
  };
in
{
  # pkgs ? import <nixpkgs> { }
  pkgs ? import nixpkgs {
    config.packageOverrides = pkgs: {
      nur = import nur {
        inherit pkgs;
      };
    };
  },
  unstable ? import nixpkgs-unstable {}
}: let
  commitlint = pkgs.buildGoModule rec {
    pname = "commitlint";
    version = "0.10.1";

    src = pkgs.fetchFromGitHub {
      owner = "conventionalcommit";
      repo = "commitlint";
      rev = "v${version}";
      hash = "sha256-OJCK6GEfs/pcorIcKjylBhdMt+lAzsBgBVUmdLfcJR0=";
    };

    # vendorHash = pkgs.lib.fakeHash;
    vendorHash = "sha256-4fV75e1Wqxsib0g31+scwM4DYuOOrHpRgavCOGurjT8=";

    meta = with pkgs.lib; {
      description = "commitlint checks if your commit messages meets the conventional commit format";
      homepage = "https://github.com/conventionalcommit/commitlint";
      license = licenses.mit;
      maintainers = with maintainers; [ wwmoraes ];
    };
  };
  # container-structure-test = pkgs.buildGoModule rec {
  #   pname = "container-structure-test";
  #   version = "1.18.1";

  #   src = pkgs.fetchFromGitHub {
  #     owner = "GoogleContainerTools";
  #     repo = "container-structure-test";
  #     rev = "v${version}";
  #     hash = "sha256-k6KP6XWeYqQhBL4Qc1CLntNCIcxS3VmGDaATCQBHO3E=";
  #   };

  #   vendorHash = "sha256-PplWNH4mtc3Vx6aGWQvUI6rcxbaTi/ovWGYDPsTyUXw=";

  #   doCheck = false;

  #   meta = with pkgs.lib; {
  #     description = "validate the structure of your container images";
  #     license = licenses.asl20;
  #     maintainers = with maintainers; [ wwmoraes ];
  #   };
  # };
  # golangci-lint = pkgs.buildGoModule rec {
  #   pname = "golangci-lint";
  #   version = "1.59.1";

  #   src = pkgs.fetchFromGitHub {
  #     owner = "golangci";
  #     repo = "golangci-lint";
  #     rev = "v${version}";
  #     hash = "sha256-VFU/qGyKBMYr0wtHXyaMjS5fXKAHWe99wDZuSyH8opg=";
  #   };

  #   vendorHash = "sha256-yYwYISK1wM/mSlAcDSIwYRo8sRWgw2u+SsvgjH+Z/7M=";

  #   subPackages = [
  #     "cmd/golangci-lint"
  #   ];

  #   ldflags = [
  #     "-s"
  #     "-w"
  #     "-X main.version=${version}"
  #     "-X main.commit=v${version}"
  #     "-X main.date=1970-01-01T00:00:00Z"
  #   ];

  #   CGO_ENABLED = 0;

  #   meta = with pkgs.lib; {
  #     description = "Fast linters runner for Go";
  #     homepage = "https://github.com/golangci/golangci-lint";
  #     license = licenses.gpl3Plus;
  #     maintainers = with maintainers; [ wwmoraes ];
  #   };
  # };
  # hadolint-sarif = pkgs.rustPlatform.buildRustPackage rec {
  #   pname = "hadolint-sarif";
  #   version = "0.4.2";

  #   src = pkgs.fetchFromGitHub {
  #     owner = "psastras";
  #     repo = "sarif-rs";
  #     rev = "hadolint-sarif-v${version}";
  #     hash = "sha256-EzWzDeIeSJ11CVcVyAhMjYQJcKHnieRrFkULc5eXAno=";
  #   };

  #   cargoHash = "sha256-AMRL1XANyze8bJe3fdgZvBnl/NyuWP13jixixqiPmiw=";
  #   cargoBuildFlags = [
  #     "--package"
  #     pname
  #   ];
  #   cargoTestFlags = cargoBuildFlags;

  #   doCheck = false;

  #   meta = with pkgs.lib; {
  #     description = "CLI tool to convert hadolint diagnostics into SARIF.";
  #     homepage = "https://github.com/psastras/sarif-rs";
  #     license = licenses.mit;
  #     maintainers = [];
  #   };
  # };
  # sarif-fmt = pkgs.rustPlatform.buildRustPackage rec {
  #   pname = "sarif-fmt";
  #   version = "0.4.2";

  #   src = pkgs.fetchFromGitHub {
  #     owner = "psastras";
  #     repo = "sarif-rs";
  #     rev = "sarif-fmt-v${version}";
  #     hash = "sha256-EzWzDeIeSJ11CVcVyAhMjYQJcKHnieRrFkULc5eXAno=";
  #   };

  #   cargoHash = "sha256-dHOxVLXtnqSHMX5r1wFxqogDf9QdnOZOjTyYFahru34=";
  #   cargoBuildFlags = [
  #     "--package"
  #     pname
  #   ];
  #   cargoTestFlags = cargoBuildFlags;

  #   doCheck = false;

  #   meta = with pkgs.lib; {
  #     description = "CLI tool to pretty print SARIF diagnostics.";
  #     homepage = "https://github.com/psastras/sarif-rs";
  #     license = licenses.mit;
  #     maintainers = [];
  #   };
  # };
in with pkgs; mkShell {
  packages = [
    commitlint
    curl
    docker-client
    editorconfig-checker
    git
    go-task
    goreleaser
    grype
    hadolint
    jq
    lefthook
    markdownlint-cli
    oapi-codegen
    python312Packages.codecov
    svu
    typos
    unstable.container-structure-test
    unstable.go
    unstable.golangci-lint
    unstable.hadolint-sarif
    unstable.sarif-fmt
    ## TODO github.com/wadey/gocovmerge
    ## TODO github.com/Khan/genqlient
    ## TODO github.com/xo/xo
  ] ++ pkgs.lib.optionals (builtins.getEnv "CI" != "") [ # CI-only
  ] ++ pkgs.lib.optionals (builtins.getEnv "CI" == "") [ # local-only
    # fish
    flyctl
    plantuml
    ## TODO pkgsite
    redis
    ## TODO structurizr-cli
  ];

  # installPhase = ''
  #   source $stdenv/setup
  # '';
}
