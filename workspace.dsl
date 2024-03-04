workspace {

  model {
    properties {
      "structurizr.groupSeparator" "/"
    }

    softwareSystem Anilist {
      anilistGraphql = container "GraphQL API"
    }

    softwareSystem Github {
      github = container "GitHub"
    }

    softwareSystem Anilistarr "converts anime sources for *arr services" {
      redis = container Redis
      bolt = container BoltDB
      badger = container BadgerDB
      postreSQL = container PostgreSQL

      container anilistarr {
        group "drivers" {
          group "trackers" {
            anilistDriver = component anilist {
              technology "struct"
              -> anilistGraphql "Uses" "HTTP/GraphQL"
            }
          }

          group "persistence" {
            redisDriver = component "Redis Driver" {
              -> redis "uses"
            }

            boltDriver = component "BoltDB Driver" {
              -> bolt "uses"
            }

            badgerDriver = component "BadgerDB Driver" {
              technology "struct"
              -> badger "Uses"
            }

            sql = component "SQL Driver" {
              technology "struct"
              -> postreSQL "Uses"
            }
          }

          group "providers" {
            fribbsProvider = component "Anilist Fribbs Provider" {
              technology "struct"
              -> github "Uses" "HTTP"
            }
          }
        }

        group "adapters" {
          group "mapper" {
            jsonProvider = component "JSON Provider" {
              technology "struct"
              -> fribbsProvider "Inherited by"
            }

            provider = component Provider {
              technology "interface"
              -> jsonProvider "Implemented by"
            }

            store = component Store {
              technology "interface"
              -> badgerDriver "Implemented by"
              -> sql "Implemented by"
            }

            mapperAdapter = component TrackerMapper {
              -> store "Uses"
              -> provider "Uses"
            }
          }

          group "cache" {
            cache = component "Cache" {
              technology "interface"
              -> redisDriver "Implemented by"
              -> boltDriver "Implemented by"
              -> badgerDriver "Implemented by"
            }

            cachedTracker = component CachedTracker {
              -> cache "Uses"
            }
          }
        }

        group "usecases" {
          tracker = component Tracker {
            technology "interface"
            -> anilistDriver "Implemented by"
            -> cachedTracker "Implemented by"
            cachedTracker -> this "Uses"
          }

          mapper = component Mapper {
            technology "interface"
            -> mapperAdapter "Implemented by"
          }

          mediaLister = component MediaLister {
            technology "struct"
            -> tracker "Uses"
            -> mapper "Uses"
          }
        }

        restAPI = component "REST API" {
          -> mediaLister "Uses"
        }
      }
    }

    softwareSystem Sonarr {
      -> restAPI "Uses" "HTTP/REST"
    }

    softwareSystem Radarr {
      -> restAPI "Uses" "HTTP/REST"
    }
  }

  views {

  }
}
