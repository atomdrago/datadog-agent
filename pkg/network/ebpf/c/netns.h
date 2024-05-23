#ifndef __NETNS_H
#define __NETNS_H

#ifndef COMPILE_CORE
#include "kconfig.h"
#include <net/net_namespace.h>  // IWYU pragma: keep // for possible_net_t and net
#include <net/sock.h>           // for sock
#include <linux/ns_common.h>    // IWYU pragma: keep // for _LINUX_NS_COMMON_H, ns_common
#endif

#include "bpf_helpers.h"        // for NULL, __always_inline
#include "bpf_telemetry.h"      // for FN_INDX_bpf_probe_read_kernel, bpf_probe_read_kernel_with_telemetry
#include "bpf_tracing.h"        // for BPF_CORE_READ_INTO, bpf_core_field_exists, BPF_PROBE_READ_INTO
#include "ktypes.h"             // for u32, __u32
#include "offsets.h"            // for offset_netns, offset_ino

#ifdef COMPILE_PREBUILT

static __always_inline __maybe_unused __u32 get_netns_from_sock(struct sock* sk) {
    void* skc_net = NULL;
    __u32 net_ns_inum = 0;
    bpf_probe_read_kernel_with_telemetry(&skc_net, sizeof(void*), ((char*)sk) + offset_netns());
    bpf_probe_read_kernel_with_telemetry(&net_ns_inum, sizeof(net_ns_inum), ((char*)skc_net) + offset_ino());
    return net_ns_inum;
}

#endif // COMPILE_PREBUILT


#ifdef COMPILE_CORE
struct sock_common___old; // iwyu insists on this forward declaration, despite it being declared below

#define sk_net __sk_common.skc_net

struct nf_conn___old {
    struct net *ct_net;
};

struct net___old {
    unsigned int proc_inum;
};

struct sock_common___old {
    struct net *skc_net;
};

struct sock___old {
    struct sock_common___old __sk_common;
};

static __always_inline __maybe_unused __u32 get_netns_from_sock(struct sock* sk) {
    u32 net_ns_inum = 0;
    struct net *ns = NULL;
    if (bpf_core_field_exists(sk->sk_net.net) ||
        bpf_core_field_exists(((struct sock___old*)sk)->sk_net->ns)) {
        BPF_CORE_READ_INTO(&ns, sk, sk_net);
        BPF_CORE_READ_INTO(&net_ns_inum, ns, ns.inum);
    } else if (bpf_core_field_exists(((struct net___old*)ns)->proc_inum)) {
        BPF_CORE_READ_INTO(&ns, (struct sock___old*)sk, sk_net);
        BPF_CORE_READ_INTO(&net_ns_inum, (struct net___old*)ns, proc_inum);
    }
    return net_ns_inum;
}

#endif // COMPILE_CORE


#ifdef COMPILE_RUNTIME

static __always_inline __maybe_unused u32 get_netns_from_sock(struct sock *sk) {
    // Retrieve network namespace id
    //
    // `possible_net_t skc_net`
    // replaced
    // `struct net *skc_net`
    // https://github.com/torvalds/linux/commit/0c5c9fb55106333e773de8c9dd321fa8240caeb3
    u32 net_ns_inum = 0;
#ifdef CONFIG_NET_NS
    struct net *ns = NULL;
    BPF_PROBE_READ_INTO(&ns, sk, sk_net);
#ifdef _LINUX_NS_COMMON_H // from: include/linux/ns_common.h
    BPF_PROBE_READ_INTO(&net_ns_inum, ns, ns.inum);
#else
    BPF_PROBE_READ_INTO(&net_ns_inum, ns, proc_inum);
#endif // LINUX_NS_COMMON_H
#endif // CONFIG_NET_NS

    return net_ns_inum;
}

#endif // COMPILE_RUNTIME

#endif
