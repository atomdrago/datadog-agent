#ifndef __BPF_TRACING_CUSTOM_H__
#define __BPF_TRACING_CUSTOM_H__

#if defined(bpf_target_x86)

#define __PT_PARM6_REG r9
#define PT_REGS_STACK_PARM(x, n)                                             \
    ({                                                                       \
        __u64 p = 0;                                                         \
        bpf_probe_read_kernel(&p, sizeof(p), ((__u64 *)x->__PT_SP_REG) + n); \
        p;                                                                   \
    })

#define PT_REGS_USER_STACK_PARM(x, n, ret)                                       \
    ({                                                                           \
        __u64 p = 0;                                                             \
        ret = bpf_probe_read_user(&p, sizeof(p), ((__u64 *)x->__PT_SP_REG) + n); \
        p;                                                                       \
    })

#define PT_REGS_PARM7(x) PT_REGS_STACK_PARM(x, 1)
#define PT_REGS_PARM8(x) PT_REGS_STACK_PARM(x, 2)
#define PT_REGS_PARM9(x) PT_REGS_STACK_PARM(x, 3)
#define PT_REGS_PARM10(x) PT_REGS_STACK_PARM(x, 4)

#define PT_REGS_USER_PARM7(x, ret) PT_REGS_USER_STACK_PARM(x, 1, ret)
#define PT_REGS_USER_PARM8(x, ret) PT_REGS_USER_STACK_PARM(x, 2, ret)
#define PT_REGS_USER_PARM9(x, ret) PT_REGS_USER_STACK_PARM(x, 3, ret)
#define PT_REGS_USER_PARM10(x, ret) PT_REGS_USER_STACK_PARM(x, 4, ret)

#elif defined(bpf_target_arm64)

#define __PT_PARM6_REG regs[5]
#define PT_REGS_STACK_PARM(x, n)                                            \
    ({                                                                      \
        unsigned long p = 0;                                                \
        bpf_probe_read_kernel(&p, sizeof(p), ((unsigned long *)x->sp) + n); \
        p;                                                                  \
    })

#define PT_REGS_USER_STACK_PARM(x, n, ret)                                      \
    ({                                                                          \
        unsigned long p = 0;                                                    \
        ret = bpf_probe_read_user(&p, sizeof(p), ((unsigned long *)x->sp) + n); \
        p;                                                                      \
    })

#define PT_REGS_PARM7(x) (__PT_REGS_CAST(x)->regs[6])
#define PT_REGS_PARM8(x) (__PT_REGS_CAST(x)->regs[7])
#define PT_REGS_PARM9(x) PT_REGS_STACK_PARM(__PT_REGS_CAST(x), 0)
#define PT_REGS_PARM10(x) PT_REGS_STACK_PARM(__PT_REGS_CAST(x), 1)
#define PT_REGS_PARM7_CORE(x) BPF_CORE_READ(__PT_REGS_CAST(x), regs[6])
#define PT_REGS_PARM8_CORE(x) BPF_CORE_READ(__PT_REGS_CAST(x), regs[7])

#define PT_REGS_USER_PARM7(x, ret) ({ \
    ret = 0;                          \
    PT_REGS_PARM7(x);                 \
})
#define PT_REGS_USER_PARM8(x, ret) ({ \
    ret = 0;                          \
    PT_REGS_PARM8(x);                 \
})

#define PT_REGS_USER_PARM9(x, ret) PT_REGS_USER_STACK_PARM(__PT_REGS_CAST(x), 0, ret)
#define PT_REGS_USER_PARM10(x, ret) PT_REGS_USER_STACK_PARM(__PT_REGS_CAST(x), 1, ret)

#endif /* defined(bpf_target_x86) */

#if defined(bpf_target_defined)

#define PT_REGS_PARM6(x) (__PT_REGS_CAST(x)->__PT_PARM6_REG)
#define PT_REGS_PARM6_CORE(x) BPF_CORE_READ(__PT_REGS_CAST(x), __PT_PARM6_REG)

#else /* defined(bpf_target_defined) */

#define PT_REGS_PARM6(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })
#define PT_REGS_PARM7(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })
#define PT_REGS_PARM8(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })
#define PT_REGS_PARM9(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })
#define PT_REGS_PARM6_CORE(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })
#define PT_REGS_PARM7_CORE(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })
#define PT_REGS_PARM8_CORE(x) ({ _Pragma(__BPF_TARGET_MISSING); 0l; })

#endif

#define ___bpf_kprobe_args6(x, args...) ___bpf_kprobe_args5(args), (void *)PT_REGS_PARM6(ctx)
#define ___bpf_kprobe_args7(x, args...) ___bpf_kprobe_args6(args), (void *)PT_REGS_PARM7(ctx)
#define ___bpf_kprobe_args8(x, args...) ___bpf_kprobe_args7(args), (void *)PT_REGS_PARM8(ctx)
#define ___bpf_kprobe_args9(x, args...) ___bpf_kprobe_args8(args), (void *)PT_REGS_PARM9(ctx)

#endif
