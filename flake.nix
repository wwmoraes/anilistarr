{
  description = "anilist custom list provider for sonarr/radarr";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-utils = {
      inputs.systems.follows = "systems";
      url = "github:numtide/flake-utils";
    };
    gomod2nix = {
      inputs.flake-utils.follows = "flake-utils";
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:tweag/gomod2nix";
    };
    nixpkgs.url = "github:NixOS/nixpkgs/25.11";
    nur = {
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-parts.follows = "flake-parts";
      url = "github:nix-community/NUR";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:numtide/treefmt-nix";
    };
    unstable.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  nixConfig = {
    substituters = [
      "https://cache.nixos.org/"
      "https://nix-community.cachix.org/"
      "https://wwmoraes.cachix.org/"
    ];
    trusted-public-keys = [
      "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
      "wwmoraes.cachix.org-1:N38Kgu19R66Jr62aX5rS466waVzT5p/Paq1g6uFFVyM="
    ];
  };

  outputs =
    inputs@{
      self,
      ...
    }:
    (inputs.flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.treefmt-nix.flakeModule
      ];

      flake = {
        overlays = {
          default = final: prev: {
            anilistarr = self.packages.${prev.stdenv.hostPlatform.system}.handler;
          };
        };
      };

      perSystem =
        {
          pkgs,
          self',
          system,
          ...
        }:
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = [
              inputs.gomod2nix.overlays.default
              inputs.nur.overlays.default
              self.overlays.default
              (final: prev: {
                unstable = import inputs.unstable { inherit (prev.stdenv.hostPlatform) system; };
              })
            ];
            config = { };
          };

          devShells = import ./shell.nix { inherit pkgs; };

          packages = {
            default = self'.packages.handler;
            handler = pkgs.callPackage ./default.nix { inherit pkgs; };
          };
        };

      systems = import inputs.systems;
    });
}
