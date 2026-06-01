import type { Directive, DirectiveBinding } from "vue";
import { useUserStore } from "@/stores/user";

export const permission: Directive = {
  mounted(el: HTMLElement, binding: DirectiveBinding) {
    const userStore = useUserStore();
    const { value } = binding;

    if (value) {
      const hasPerm = Array.isArray(value)
        ? value.some((p: string) => userStore.hasPermission(p))
        : userStore.hasPermission(value);

      if (!hasPerm) {
        el.parentNode?.removeChild(el);
      }
    }
  },
};
