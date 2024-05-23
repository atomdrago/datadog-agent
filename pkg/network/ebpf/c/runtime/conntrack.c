#ifdef COMPILE_RUNTIME
#include "kconfig.h"
#include <net/netfilter/nf_conntrack.h>                 // for nf_conn
#include <uapi/linux/netfilter/nf_conntrack_common.h>   // for IPS_CONFIRMED, IPS_NAT_MASK
#endif

#include "conntrack.h"          // for nf_conn_to_conntrack_tuples
#include "bpf_helpers.h"        // for SEC, log_debug, bpf_get_current_pid_tgid
#include "bpf_telemetry.h"      // for bpf_map_update_with_telemetry
#include "bpf_tracing.h"        // for BPF_KPROBE, PT_REGS_PARM5, BPF_CORE_READ_INTO
#include "conntrack/helpers.h"  // for increment_telemetry_registers_count
#include "conntrack/maps.h"     // for conntrack
#include "ktypes.h"             // for u32, BPF_ANY, pt_regs

SEC("kprobe/__nf_conntrack_hash_insert")
int BPF_KPROBE(kprobe___nf_conntrack_hash_insert, struct nf_conn *ct) {
    u32 status = 0;
    BPF_CORE_READ_INTO(&status, ct, status);
    if (!(status&IPS_CONFIRMED) || !(status&IPS_NAT_MASK)) {
        return 0;
    }

    log_debug("kprobe/__nf_conntrack_hash_insert: netns: %u, status: %x", get_netns(ct), status);

    conntrack_tuple_t orig = {}, reply = {};
    if (nf_conn_to_conntrack_tuples(ct, &orig, &reply) != 0) {
        return 0;
    }

    bpf_map_update_with_telemetry(conntrack, &orig, &reply, BPF_ANY);
    bpf_map_update_with_telemetry(conntrack, &reply, &orig, BPF_ANY);
    increment_telemetry_registers_count();

    return 0;
}

SEC("kprobe/ctnetlink_fill_info")
int BPF_KPROBE(kprobe_ctnetlink_fill_info) {
    u32 pid = bpf_get_current_pid_tgid() >> 32;
    if (pid != systemprobe_pid()) {
        log_debug("skipping kprobe/ctnetlink_fill_info invocation from non-system-probe process");
        return 0;
    }

    struct nf_conn *ct = (struct nf_conn *)PT_REGS_PARM5(ctx);

    u32 status = 0;
    BPF_CORE_READ_INTO(&status, ct, status);
    if (!(status&IPS_CONFIRMED) || !(status&IPS_NAT_MASK)) {
        return 0;
    }

    log_debug("kprobe/ctnetlink_fill_info: netns: %u, status: %x", get_netns(ct), status);

    conntrack_tuple_t orig = {}, reply = {};
    if (nf_conn_to_conntrack_tuples(ct, &orig, &reply) != 0) {
        return 0;
    }

    bpf_map_update_with_telemetry(conntrack, &orig, &reply, BPF_ANY);
    bpf_map_update_with_telemetry(conntrack, &reply, &orig, BPF_ANY);
    increment_telemetry_registers_count();

    return 0;
}

char _license[] SEC("license") = "GPL";
