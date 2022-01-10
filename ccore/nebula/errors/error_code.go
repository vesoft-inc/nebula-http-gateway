package errors

import (
	stderrors "errors"
	"fmt"
)

type (
	ErrorCode int64

	codeError struct {
		errorCode ErrorCode
		errorMsg  string
	}

	CodeError interface {
		error
		GetErrorCode() ErrorCode
		GetErrorMsg() string
	}
)

func AsCodeError(err error) (CodeError, bool) {
	if e := new(codeError); stderrors.As(err, &e) {
		return e, true
	}
	return nil, false
}

func IsCodeError(err error, c ...ErrorCode) bool {
	ce, ok := AsCodeError(err)
	if !ok {
		return false
	}
	switch len(c) {
	case 0:
		return true
	case 1:
		return ce.GetErrorCode() == c[0]
	default:
		return false
	}
}

func NewCodeError(errorCode ErrorCode, errorMsg string) error {
	return &codeError{
		errorCode: errorCode,
		errorMsg:  errorMsg,
	}
}

func (e *codeError) Error() string {
	return fmt.Sprintf("%d:%s", int64(e.errorCode), e.errorMsg)
}

func (e *codeError) GetErrorCode() ErrorCode {
	return e.errorCode
}

func (e *codeError) GetErrorMsg() string {
	return e.errorMsg
}

const (
	ErrorCode_SUCCEEDED                           ErrorCode = 0
	ErrorCode_E_DISCONNECTED                      ErrorCode = -1
	ErrorCode_E_FAIL_TO_CONNECT                   ErrorCode = -2
	ErrorCode_E_RPC_FAILURE                       ErrorCode = -3
	ErrorCode_E_LEADER_CHANGED                    ErrorCode = -4
	ErrorCode_E_SPACE_NOT_FOUND                   ErrorCode = -5
	ErrorCode_E_TAG_NOT_FOUND                     ErrorCode = -6
	ErrorCode_E_EDGE_NOT_FOUND                    ErrorCode = -7
	ErrorCode_E_INDEX_NOT_FOUND                   ErrorCode = -8
	ErrorCode_E_EDGE_PROP_NOT_FOUND               ErrorCode = -9
	ErrorCode_E_TAG_PROP_NOT_FOUND                ErrorCode = -10
	ErrorCode_E_ROLE_NOT_FOUND                    ErrorCode = -11
	ErrorCode_E_CONFIG_NOT_FOUND                  ErrorCode = -12
	ErrorCode_E_GROUP_NOT_FOUND                   ErrorCode = -13
	ErrorCode_E_ZONE_NOT_FOUND                    ErrorCode = -14
	ErrorCode_E_LISTENER_NOT_FOUND                ErrorCode = -15
	ErrorCode_E_PART_NOT_FOUND                    ErrorCode = -16
	ErrorCode_E_KEY_NOT_FOUND                     ErrorCode = -17
	ErrorCode_E_USER_NOT_FOUND                    ErrorCode = -18
	ErrorCode_E_STATS_NOT_FOUND                   ErrorCode = -19
	ErrorCode_E_BACKUP_FAILED                     ErrorCode = -24
	ErrorCode_E_BACKUP_EMPTY_TABLE                ErrorCode = -25
	ErrorCode_E_BACKUP_TABLE_FAILED               ErrorCode = -26
	ErrorCode_E_PARTIAL_RESULT                    ErrorCode = -27
	ErrorCode_E_REBUILD_INDEX_FAILED              ErrorCode = -28
	ErrorCode_E_INVALID_PASSWORD                  ErrorCode = -29
	ErrorCode_E_FAILED_GET_ABS_PATH               ErrorCode = -30
	ErrorCode_E_BAD_USERNAME_PASSWORD             ErrorCode = -1001
	ErrorCode_E_SESSION_INVALID                   ErrorCode = -1002
	ErrorCode_E_SESSION_TIMEOUT                   ErrorCode = -1003
	ErrorCode_E_SYNTAX_ERROR                      ErrorCode = -1004
	ErrorCode_E_EXECUTION_ERROR                   ErrorCode = -1005
	ErrorCode_E_STATEMENT_EMPTY                   ErrorCode = -1006
	ErrorCode_E_BAD_PERMISSION                    ErrorCode = -1008
	ErrorCode_E_SEMANTIC_ERROR                    ErrorCode = -1009
	ErrorCode_E_TOO_MANY_CONNECTIONS              ErrorCode = -1010
	ErrorCode_E_PARTIAL_SUCCEEDED                 ErrorCode = -1011
	ErrorCode_E_NO_HOSTS                          ErrorCode = -2001
	ErrorCode_E_EXISTED                           ErrorCode = -2002
	ErrorCode_E_INVALID_HOST                      ErrorCode = -2003
	ErrorCode_E_UNSUPPORTED                       ErrorCode = -2004
	ErrorCode_E_NOT_DROP                          ErrorCode = -2005
	ErrorCode_E_BALANCER_RUNNING                  ErrorCode = -2006
	ErrorCode_E_CONFIG_IMMUTABLE                  ErrorCode = -2007
	ErrorCode_E_CONFLICT                          ErrorCode = -2008
	ErrorCode_E_INVALID_PARM                      ErrorCode = -2009
	ErrorCode_E_WRONGCLUSTER                      ErrorCode = -2010
	ErrorCode_E_STORE_FAILURE                     ErrorCode = -2021
	ErrorCode_E_STORE_SEGMENT_ILLEGAL             ErrorCode = -2022
	ErrorCode_E_BAD_BALANCE_PLAN                  ErrorCode = -2023
	ErrorCode_E_BALANCED                          ErrorCode = -2024
	ErrorCode_E_NO_RUNNING_BALANCE_PLAN           ErrorCode = -2025
	ErrorCode_E_NO_VALID_HOST                     ErrorCode = -2026
	ErrorCode_E_CORRUPTTED_BALANCE_PLAN           ErrorCode = -2027
	ErrorCode_E_NO_INVALID_BALANCE_PLAN           ErrorCode = -2028
	ErrorCode_E_IMPROPER_ROLE                     ErrorCode = -2030
	ErrorCode_E_INVALID_PARTITION_NUM             ErrorCode = -2031
	ErrorCode_E_INVALID_REPLICA_FACTOR            ErrorCode = -2032
	ErrorCode_E_INVALID_CHARSET                   ErrorCode = -2033
	ErrorCode_E_INVALID_COLLATE                   ErrorCode = -2034
	ErrorCode_E_CHARSET_COLLATE_NOT_MATCH         ErrorCode = -2035
	ErrorCode_E_SNAPSHOT_FAILURE                  ErrorCode = -2040
	ErrorCode_E_BLOCK_WRITE_FAILURE               ErrorCode = -2041
	ErrorCode_E_REBUILD_INDEX_FAILURE             ErrorCode = -2042
	ErrorCode_E_INDEX_WITH_TTL                    ErrorCode = -2043
	ErrorCode_E_ADD_JOB_FAILURE                   ErrorCode = -2044
	ErrorCode_E_STOP_JOB_FAILURE                  ErrorCode = -2045
	ErrorCode_E_SAVE_JOB_FAILURE                  ErrorCode = -2046
	ErrorCode_E_BALANCER_FAILURE                  ErrorCode = -2047
	ErrorCode_E_JOB_NOT_FINISHED                  ErrorCode = -2048
	ErrorCode_E_TASK_REPORT_OUT_DATE              ErrorCode = -2049
	ErrorCode_E_JOB_NOT_IN_SPACE                  ErrorCode = -2050
	ErrorCode_E_INVALID_JOB                       ErrorCode = -2065
	ErrorCode_E_BACKUP_BUILDING_INDEX             ErrorCode = -2066
	ErrorCode_E_BACKUP_SPACE_NOT_FOUND            ErrorCode = -2067
	ErrorCode_E_RESTORE_FAILURE                   ErrorCode = -2068
	ErrorCode_E_SESSION_NOT_FOUND                 ErrorCode = -2069
	ErrorCode_E_LIST_CLUSTER_FAILURE              ErrorCode = -2070
	ErrorCode_E_LIST_CLUSTER_GET_ABS_PATH_FAILURE ErrorCode = -2071
	ErrorCode_E_GET_META_DIR_FAILURE              ErrorCode = -2072
	ErrorCode_E_QUERY_NOT_FOUND                   ErrorCode = -2073
	ErrorCode_E_CONSENSUS_ERROR                   ErrorCode = -3001
	ErrorCode_E_KEY_HAS_EXISTS                    ErrorCode = -3002
	ErrorCode_E_DATA_TYPE_MISMATCH                ErrorCode = -3003
	ErrorCode_E_INVALID_FIELD_VALUE               ErrorCode = -3004
	ErrorCode_E_INVALID_OPERATION                 ErrorCode = -3005
	ErrorCode_E_NOT_NULLABLE                      ErrorCode = -3006
	ErrorCode_E_FIELD_UNSET                       ErrorCode = -3007
	ErrorCode_E_OUT_OF_RANGE                      ErrorCode = -3008
	ErrorCode_E_ATOMIC_OP_FAILED                  ErrorCode = -3009
	ErrorCode_E_DATA_CONFLICT_ERROR               ErrorCode = -3010
	ErrorCode_E_WRITE_STALLED                     ErrorCode = -3011
	ErrorCode_E_IMPROPER_DATA_TYPE                ErrorCode = -3021
	ErrorCode_E_INVALID_SPACEVIDLEN               ErrorCode = -3022
	ErrorCode_E_INVALID_FILTER                    ErrorCode = -3031
	ErrorCode_E_INVALID_UPDATER                   ErrorCode = -3032
	ErrorCode_E_INVALID_STORE                     ErrorCode = -3033
	ErrorCode_E_INVALID_PEER                      ErrorCode = -3034
	ErrorCode_E_RETRY_EXHAUSTED                   ErrorCode = -3035
	ErrorCode_E_TRANSFER_LEADER_FAILED            ErrorCode = -3036
	ErrorCode_E_INVALID_STAT_TYPE                 ErrorCode = -3037
	ErrorCode_E_INVALID_VID                       ErrorCode = -3038
	ErrorCode_E_NO_TRANSFORMED                    ErrorCode = -3039
	ErrorCode_E_LOAD_META_FAILED                  ErrorCode = -3040
	ErrorCode_E_FAILED_TO_CHECKPOINT              ErrorCode = -3041
	ErrorCode_E_CHECKPOINT_BLOCKED                ErrorCode = -3042
	ErrorCode_E_FILTER_OUT                        ErrorCode = -3043
	ErrorCode_E_INVALID_DATA                      ErrorCode = -3044
	ErrorCode_E_MUTATE_EDGE_CONFLICT              ErrorCode = -3045
	ErrorCode_E_MUTATE_TAG_CONFLICT               ErrorCode = -3046
	ErrorCode_E_OUTDATED_LOCK                     ErrorCode = -3047
	ErrorCode_E_INVALID_TASK_PARA                 ErrorCode = -3051
	ErrorCode_E_USER_CANCEL                       ErrorCode = -3052
	ErrorCode_E_TASK_EXECUTION_FAILED             ErrorCode = -3053
	ErrorCode_E_PLAN_IS_KILLED                    ErrorCode = -3060
	ErrorCode_E_NO_TERM                           ErrorCode = -3070
	ErrorCode_E_OUTDATED_TERM                     ErrorCode = -3071
	ErrorCode_E_OUTDATED_EDGE                     ErrorCode = -3072
	ErrorCode_E_WRITE_WRITE_CONFLICT              ErrorCode = -3073
	ErrorCode_E_CLIENT_SERVER_INCOMPATIBLE        ErrorCode = -3061
	ErrorCode_E_UNKNOWN                           ErrorCode = -8000
)
