workspace {

  model {
    properties {
      "structurizr.groupSeparator" "/"
    }

    softwareSystem Anilist {
      anilistAPI = container API {
        technology "GraphQL"
        tags "API"
      }
    }

    softwareSystem GitHub {
      githubAPI = container API {
        technology "REST"
        tags "API"
      }
    }

    softwareSystem Anilistarr "converts anime sources for *arr services" {
      redis = container Redis {
        tags "Database"
      }

      bolt = container BoltDB {
        tags "Database"
      }

      badger = container BadgerDB {
        tags "Database"
      }

      sqlite = container SQLite {
        tags "Database"
      }

      container anilistarr {
        group "drivers" {
          group "trackers" {
            anilistDriver = component "Anilist Driver" {
              technology "struct"
              -> anilistAPI "Uses" "HTTP/GraphQL"
            }
          }

          group "caches | stores" {
            badgerDriver = component "BadgerDB Driver" {
              -> badger "Uses" "IPC"
            }

            boltDriver = component "BoltDB Driver" {
              -> bolt "uses" "IPC"
            }

            redisDriver = component "Redis Driver" {
              -> redis "uses" "HTTP"
            }

            sqliteDriver = component "SQLite Driver" {
              -> SQLite "Uses" "IPC"
            }
          }

          group "providers" {
            fribbsProvider = component "Anilist Fribbs Provider" {
              -> githubAPI "Uses" "HTTP/REST"
            }
          }
        }

        group "usecases" {
          cachedTracker = component CachedTracker {
            -> anilistDriver "Uses"

            -> boltDriver  "Uses"
            -> redisDriver  "Uses"
            -> badgerDriver  "Uses"
            -> sqliteDriver  "Uses"
          }

          mediaList = component MediaList {
            -> cachedTracker "Uses"
            -> fribbsProvider "Uses"
            -> badgerDriver "Uses"
            -> sqliteDriver "Uses"
          }
        }

        restAPI = component "API" {
          technology "REST"
          tags "API"

          -> mediaList "Uses"
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
    styles {
      element "Person" {
        shape Person
      }
      element "API" {
        shape hexagon
      }
      element "Database" {
        shape cylinder
      }
    }
  }
}
