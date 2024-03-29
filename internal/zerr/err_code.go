package zerr

const (
	REQUEST_OK          int = 0x0000
	MYSQL_CONNECT_ERROR int = 0x0001
	MYSQL_CLOSE_ERROR   int = 0x0002
	MYSQL_ERROR   int = 0x0003
	MYSQL_NO_RESULT     int = 0x0004


	REQUEST_PARAM_ERROR int = 0x0010
	REQUEST_ERROR int = 0x0011


	SERVER_ERROR        int = 0x0020


	REMOTE_GRPC_CONNECT_ERROR   int = 0x0030
	REMOTE_GRPC_TIMEOUT_ERROR   int = 0x0031
	REMOTE_GRPC_ERROR   int = 0x0032

	REDIS_CONNECT_ERROR   int = 0x0040
	REDIS_CLOSE_ERROR   int = 0x0041
	REDIS_ERROR   int = 0x0041

	ETCD_PARAM_ERROR int = 0x0050
	ETCD_READ_ERROR int = 0x0051
	ETCD_WRITE_ERROR int = 0x0052
	ETCD_DELETE_ERROR int = 0x0053
	ETCD_HEART_ERROR int = 0x0054
	ETCD_REGISTER_ERROR int = 0x0055
	ETCD_UNREGISTER_ERROR int = 0x0056

	LOG_WRITE_ERROR  int = 0x0060
	LOG_NO_DISK_ERROR  int = 0x0061

	CONFIG_FLAG_USAGED_ERROR  int = 0x0070
	CONFIG_FILE_LOADING_ERROR  int = 0x0071
	CONFIG_LOADING_ERROR  int = 0x0072
	CONFIG_GET_CURRENT_FILE_ERROR  int = 0x0073

	PWD_DIR_NOT_FIND_ERROR int = 0x0080
	EXECUTABLE_DIR_NOT_FIND_ERROR int = 0x0080

	UNKNOWN_ERROR   int = 0xffff
)