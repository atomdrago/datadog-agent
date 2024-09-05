// Code generated by "stringer -output fatbin_string.go -type=nvInfoAttr,nvInfoFormat,fatbinDataKind -linecomment"; DO NOT EDIT.

package gpu

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NVI_ATTR_ERROR-0]
	_ = x[NVI_ATTR_PAD-1]
	_ = x[NVI_ATTR_IMAGE_SLOT-2]
	_ = x[NVI_ATTR_JUMPTABLE_RELOCS-3]
	_ = x[NVI_ATTR_CTAIDZ_USED-4]
	_ = x[NVI_ATTR_MAX_THREADS-5]
	_ = x[NVI_ATTR_IMAGE_OFFSET-6]
	_ = x[NVI_ATTR_IMAGE_SIZE-7]
	_ = x[NVI_ATTR_TEXTURE_NORMALIZED-8]
	_ = x[NVI_ATTR_SAMPLER_INIT-9]
	_ = x[NVI_ATTR_PARAM_CBANK-10]
	_ = x[NVI_ATTR_SMEM_PARAM_OFFSETS-11]
	_ = x[NVI_ATTR_CBANK_PARAM_OFFSETS-12]
	_ = x[NVI_ATTR_SYNC_STACK-13]
	_ = x[NVI_ATTR_TEXID_SAMPID_MAP-14]
	_ = x[NVI_ATTR_EXTERNS-15]
	_ = x[NVI_ATTR_REQNTID-16]
	_ = x[NVI_ATTR_FRAME_SIZE-17]
	_ = x[NVI_ATTR_MIN_STACK_SIZE-18]
	_ = x[NVI_ATTR_SAMPLER_FORCE_UNNORMALIZED-19]
	_ = x[NVI_ATTR_BINDLESS_IMAGE_OFFSETS-20]
	_ = x[NVI_ATTR_BINDLESS_TEXTURE_BANK-21]
	_ = x[NVI_ATTR_BINDLESS_SURFACE_BANK-22]
	_ = x[NVI_ATTR_KPARAM_INFO-23]
	_ = x[NVI_ATTR_SMEM_PARAM_SIZE-24]
	_ = x[NVI_ATTR_CBANK_PARAM_SIZE-25]
	_ = x[NVI_ATTR_QUERY_NUMATTRIB-26]
	_ = x[NVI_ATTR_MAXREG_COUNT-27]
	_ = x[NVI_ATTR_EXIT_INSTR_OFFSETS-28]
	_ = x[NVI_ATTR_S2RCTAID_INSTR_OFFSETS-29]
	_ = x[NVI_ATTR_CRS_STACK_SIZE-30]
	_ = x[NVI_ATTR_NEED_CNP_WRAPPER-31]
	_ = x[NVI_ATTR_NEED_CNP_PATCH-32]
	_ = x[NVI_ATTR_EXPLICIT_CACHING-33]
	_ = x[NVI_ATTR_ISTYPEP_USED-34]
	_ = x[NVI_ATTR_MAX_STACK_SIZE-35]
	_ = x[NVI_ATTR_SUQ_USED-36]
	_ = x[NVI_ATTR_LD_CACHEMOD_INSTR_OFFSETS-37]
	_ = x[NVI_ATTR_LOAD_CACHE_REQUEST-38]
	_ = x[NVI_ATTR_ATOM_SYS_INSTR_OFFSETS-39]
	_ = x[NVI_ATTR_COOP_GROUP_INSTR_OFFSETS-40]
	_ = x[NVI_ATTR_COOP_GROUP_MAX_REGIDS-41]
	_ = x[NVI_ATTR_SW1850030_WAR-42]
	_ = x[NVI_ATTR_WMMA_USED-43]
	_ = x[NVI_ATTR_HAS_PRE_V10_OBJECT-44]
	_ = x[NVI_ATTR_ATOMF16_EMUL_INSTR_OFFSETS-45]
	_ = x[NVI_ATTR_ATOM16_EMUL_INSTR_REG_MAP-46]
	_ = x[NVI_ATTR_REGCOUNT-47]
	_ = x[NVI_ATTR_SW2393858_WAR-48]
	_ = x[NVI_ATTR_INT_WARP_WIDE_INSTR_OFFSETS-49]
	_ = x[NVI_ATTR_SHARED_SCRATCH-50]
	_ = x[NVI_ATTR_STATISTICS-51]
	_ = x[NVI_ATTR_INDIRECT_BRANCH_TARGETS-52]
	_ = x[NVI_ATTR_SW2861232_WAR-53]
	_ = x[NVI_ATTR_SW_WAR-54]
	_ = x[NVI_ATTR_CUDA_API_VERSION-55]
	_ = x[NVI_ATTR_NUM_MBARRIERS-56]
	_ = x[NVI_ATTR_MBARRIER_INSTR_OFFSETS-57]
	_ = x[NVI_ATTR_COROUTINE_RESUME_ID_OFFSETS-58]
	_ = x[NVI_ATTR_SAM_REGION_STACK_SIZE-59]
	_ = x[NVI_ATTR_PER_REG_TARGET_PERF_STATS-60]
	_ = x[NVI_ATTR_CTA_PER_CLUSTER-61]
	_ = x[NVI_ATTR_EXPLICIT_CLUSTER-62]
	_ = x[NVI_ATTR_MAX_CLUSTER_RANK-63]
	_ = x[NVI_ATTR_INSTR_REG_MAP-64]
}

const _nvInfoAttr_name = "NVI_ATTR_ERRORNVI_ATTR_PADNVI_ATTR_IMAGE_SLOTNVI_ATTR_JUMPTABLE_RELOCSNVI_ATTR_CTAIDZ_USEDNVI_ATTR_MAX_THREADSNVI_ATTR_IMAGE_OFFSETNVI_ATTR_IMAGE_SIZENVI_ATTR_TEXTURE_NORMALIZEDNVI_ATTR_SAMPLER_INITNVI_ATTR_PARAM_CBANKNVI_ATTR_SMEM_PARAM_OFFSETSNVI_ATTR_CBANK_PARAM_OFFSETSNVI_ATTR_SYNC_STACKNVI_ATTR_TEXID_SAMPID_MAPNVI_ATTR_EXTERNSNVI_ATTR_REQNTIDNVI_ATTR_FRAME_SIZENVI_ATTR_MIN_STACK_SIZENVI_ATTR_SAMPLER_FORCE_UNNORMALIZEDNVI_ATTR_BINDLESS_IMAGE_OFFSETSNVI_ATTR_BINDLESS_TEXTURE_BANKNVI_ATTR_BINDLESS_SURFACE_BANKNVI_ATTR_KPARAM_INFONVI_ATTR_SMEM_PARAM_SIZENVI_ATTR_CBANK_PARAM_SIZENVI_ATTR_QUERY_NUMATTRIBNVI_ATTR_MAXREG_COUNTNVI_ATTR_EXIT_INSTR_OFFSETSNVI_ATTR_S2RCTAID_INSTR_OFFSETSNVI_ATTR_CRS_STACK_SIZENVI_ATTR_NEED_CNP_WRAPPERNVI_ATTR_NEED_CNP_PATCHNVI_ATTR_EXPLICIT_CACHINGNVI_ATTR_ISTYPEP_USEDNVI_ATTR_MAX_STACK_SIZENVI_ATTR_SUQ_USEDNVI_ATTR_LD_CACHEMOD_INSTR_OFFSETSNVI_ATTR_LOAD_CACHE_REQUESTNVI_ATTR_ATOM_SYS_INSTR_OFFSETSNVI_ATTR_COOP_GROUP_INSTR_OFFSETSNVI_ATTR_COOP_GROUP_MAX_REGIDSNVI_ATTR_SW1850030_WARNVI_ATTR_WMMA_USEDNVI_ATTR_HAS_PRE_V10_OBJECTNVI_ATTR_ATOMF16_EMUL_INSTR_OFFSETSNVI_ATTR_ATOM16_EMUL_INSTR_REG_MAPNVI_ATTR_REGCOUNTNVI_ATTR_SW2393858_WARNVI_ATTR_INT_WARP_WIDE_INSTR_OFFSETSNVI_ATTR_SHARED_SCRATCHNVI_ATTR_STATISTICSNVI_ATTR_INDIRECT_BRANCH_TARGETSNVI_ATTR_SW2861232_WARNVI_ATTR_SW_WARNVI_ATTR_CUDA_API_VERSIONNVI_ATTR_NUM_MBARRIERSNVI_ATTR_MBARRIER_INSTR_OFFSETSNVI_ATTR_COROUTINE_RESUME_ID_OFFSETSNVI_ATTR_SAM_REGION_STACK_SIZENVI_ATTR_PER_REG_TARGET_PERF_STATSNVI_ATTR_CTA_PER_CLUSTERNVI_ATTR_EXPLICIT_CLUSTERNVI_ATTR_MAX_CLUSTER_RANKNVI_ATTR_INSTR_REG_MAP"

var _nvInfoAttr_index = [...]uint16{0, 14, 26, 45, 70, 90, 110, 131, 150, 177, 198, 218, 245, 273, 292, 317, 333, 349, 368, 391, 426, 457, 487, 517, 537, 561, 586, 610, 631, 658, 689, 712, 737, 760, 785, 806, 829, 846, 880, 907, 938, 971, 1001, 1023, 1041, 1068, 1103, 1137, 1154, 1176, 1212, 1235, 1254, 1286, 1308, 1323, 1348, 1370, 1401, 1437, 1467, 1501, 1525, 1550, 1575, 1597}

func (i nvInfoAttr) String() string {
	if i >= nvInfoAttr(len(_nvInfoAttr_index)-1) {
		return "nvInfoAttr(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _nvInfoAttr_name[_nvInfoAttr_index[i]:_nvInfoAttr_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NVI_FMT_NONE-1]
	_ = x[NVI_FMT_BVAL-2]
	_ = x[NVI_FMT_HVAL-3]
	_ = x[NVI_FMT_SVAL-4]
}

const _nvInfoFormat_name = "NVI_FMT_NONENVI_FMT_BVALNVI_FMT_HVALNVI_FMT_SVAL"

var _nvInfoFormat_index = [...]uint8{0, 12, 24, 36, 48}

func (i nvInfoFormat) String() string {
	i -= 1
	if i >= nvInfoFormat(len(_nvInfoFormat_index)-1) {
		return "nvInfoFormat(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _nvInfoFormat_name[_nvInfoFormat_index[i]:_nvInfoFormat_index[i+1]]
}
func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FATBIN_DATA_KIND_PTX-1]
	_ = x[FATBIN_DATA_KIND_SM-2]
}

const _fatbinDataKind_name = "FATBIN_DATA_KIND_PTXFATBIN_DATA_KIND_SM"

var _fatbinDataKind_index = [...]uint8{0, 20, 39}

func (i fatbinDataKind) String() string {
	i -= 1
	if i >= fatbinDataKind(len(_fatbinDataKind_index)-1) {
		return "fatbinDataKind(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _fatbinDataKind_name[_fatbinDataKind_index[i]:_fatbinDataKind_index[i+1]]
}
