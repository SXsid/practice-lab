class Database:
    def __init__(self, host, port, ssl, pool_size):
        self.__host = host
        self.__port = port
        self.__ssl = ssl
        self.__pool_size = pool_size

    def connect(self):
        pass


class DatabaseBuilder:
    def __init__(self):
        self.__host = None
        self.__port = 5432  # default
        self.__ssl = False
        self.__pool_size = 10

    def with_host(self, host: str) -> "DatabaseBuilder":
        self.__host = host
        return self

    def with_port(self, port: int) -> "DatabaseBuilder":
        self.__port = port
        return self

    def with_ssl(self, ssl: bool) -> "DatabaseBuilder":
        self.__ssl = ssl
        return self

    def with_pool_size(self, size: int) -> "DatabaseBuilder":
        self.__pool_size = size
        return self

    def build(self) -> Database:
        if not self.__host:
            raise ValueError("host is required")
        return Database(self.__host, self.__port, self.__ssl, self.__pool_size)


# usage
# INFO: main is the that how with funtion return the serlf "DatabaseBuilder"
db = DatabaseBuilder().with_host("localhost").with_ssl(True).with_pool_size(100).build()
# ```
#
# ---
#
# ## The key insight — `return self`
#
# Every `With` method returns the builder itself. That's what enables chaining:
# ```
# NewDatabase()          → returns builder
# .WithHost("localhost") → sets host, returns same builder
# .WithPort(5432)        → sets port, returns same builder
# .Build()               → validates, constructs, returns final object
