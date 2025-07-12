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
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-25.05-darwin";
    nur = {
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-parts.follows = "flake-parts";
      url = "github:nix-community/NUR";
    };
    # sops-nix = {
    #   inputs.nixpkgs.follows = "nixpkgs";
    #   url = "github:Mic92/sops-nix";
    # };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      inputs.nixpkgs.follows = "nixpkgs";
      url = "github:numtide/treefmt-nix";
    };
    unstable.url = "github:NixOS/nixpkgs?rev=e38c80c027d6bbdfa2a305fc08e732b9fac4928a";
  };

  nixConfig = {
    substituters = [
      "https://wwmoraes.cachix.org/"
      "https://nix-community.cachix.org/"
      "https://cache.nixos.org/"
    ];
    trusted-public-keys = [
      "wwmoraes.cachix.org-1:N38Kgu19R66Jr62aX5rS466waVzT5p/Paq1g6uFFVyM="
      "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
      "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="
    ];
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      gomod2nix,
      nixpkgs,
      nur,
      # sops-nix,
      systems,
      treefmt-nix,
      unstable,
      ...
    }:
    (flake-parts.lib.mkFlake { inherit inputs; } {
      flake = {
        overlays = {
          default = final: prev: {
            anilistarr = self.packages.${prev.system}.handler;
          };
        };
      };

      perSystem =
        { pkgs, system, ... }:
        let
          treefmt = treefmt-nix.lib.evalModule pkgs ./treefmt.nix;
        in
        {
          _module.args.pkgs = import nixpkgs {
            inherit system;
            overlays = [
              gomod2nix.overlays.default
              nur.overlays.default
              self.overlays.default
              (final: prev: {
                goEnv = prev.mkGoEnv { pwd = ./.; };
                unstable = import unstable { inherit (prev) system; };
              })
            ];
            config = { };
          };

          checks = {
            formatting = treefmt.config.build.check self;
          };

          formatter = treefmt.config.build.wrapper;

          devShells.default = import ./shell.nix { inherit pkgs; };

          packages = rec {
            default = handler;
            handler = pkgs.buildGoApplication {
              pname = "anilistarr";
              version = self.shortRev or self.dirtyShortRev;
              src =
                with pkgs.lib.fileset;
                toSource {
                  root = ./.;
                  fileset = intersection (gitTracked ./.) (unions [
                    (fileFilter (file: file.hasExt "go") ./.)
                    ./cmd/handler
                    ./go.mod
                    ./go.sum
                    ./internal
                    ./pkg
                    ./sqlc.yaml
                    ./swagger.yaml
                  ]);
                };
              modules = ./gomod2nix.toml;
              subPackages = [ "cmd/handler" ];
              meta = {
                description = "anilist custom list provider for sonarr/radarr";
                homepage = "https://github.com/wwmoraes/anilistarr";
                license = pkgs.lib.licenses.mit;
                maintainers = [ pkgs.lib.maintainers.wwmoraes ];
                mainProgram = "handler";
              };
            };
          };
        };

      systems = import systems;
    });
}
