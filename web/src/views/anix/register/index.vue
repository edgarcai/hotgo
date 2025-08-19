<template>
  <div
    class="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8"
  >
    <div class="max-w-md w-full space-y-8">
      <!-- 头部 -->
      <div class="text-center">
        <h2 class="mt-6 text-3xl font-extrabold text-gray-900"> 注册 AniX 账户 </h2>
        <p class="mt-2 text-sm text-gray-600">
          已有账户？
          <router-link to="/anix/login" class="font-medium text-indigo-600 hover:text-indigo-500">
            立即登录
          </router-link>
        </p>
      </div>

      <!-- 注册表单 -->
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="bg-white rounded-lg shadow-md p-6 space-y-4">
          <!-- 用户名 -->
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700"> 用户名 </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="请输入用户名"
              :class="{ 'border-red-500': errors.username }"
            />
            <p v-if="errors.username" class="mt-1 text-sm text-red-600">{{ errors.username }}</p>
          </div>

          <!-- 邮箱 -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700"> 邮箱地址 </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="请输入邮箱地址"
              :class="{ 'border-red-500': errors.email }"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
          </div>

          <!-- 昵称 -->
          <div>
            <label for="nickname" class="block text-sm font-medium text-gray-700"> 昵称 </label>
            <input
              id="nickname"
              v-model="form.nickname"
              type="text"
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="请输入昵称（可选）"
              :class="{ 'border-red-500': errors.nickname }"
            />
            <p v-if="errors.nickname" class="mt-1 text-sm text-red-600">{{ errors.nickname }}</p>
          </div>

          <!-- 密码 -->
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700"> 密码 </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="请输入密码"
              :class="{ 'border-red-500': errors.password }"
            />
            <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
            <!-- 密码强度提示 -->
            <div v-if="form.password" class="mt-2">
              <div class="flex items-center space-x-2">
                <div class="flex-1 bg-gray-200 rounded-full h-2">
                  <div
                    class="h-2 rounded-full transition-all duration-300"
                    :class="passwordStrengthClass"
                    :style="{ width: passwordStrengthWidth }"
                  ></div>
                </div>
                <span class="text-xs" :class="passwordStrengthTextClass">
                  {{ passwordStrengthText }}
                </span>
              </div>
            </div>
          </div>

          <!-- 确认密码 -->
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
              确认密码
            </label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              class="mt-1 appearance-none relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="请再次输入密码"
              :class="{ 'border-red-500': errors.confirmPassword }"
            />
            <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{
              errors.confirmPassword
            }}</p>
          </div>

          <!-- 用户协议 -->
          <div class="flex items-center">
            <input
              id="agree"
              v-model="form.agree"
              type="checkbox"
              class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            />
            <label for="agree" class="ml-2 block text-sm text-gray-900">
              我已阅读并同意
              <a href="#" class="text-indigo-600 hover:text-indigo-500">用户协议</a>
              和
              <a href="#" class="text-indigo-600 hover:text-indigo-500">隐私政策</a>
            </label>
          </div>
          <p v-if="errors.agree" class="text-sm text-red-600">{{ errors.agree }}</p>
        </div>

        <!-- 注册按钮 -->
        <div>
          <button
            type="submit"
            :disabled="loading || !form.agree"
            class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
          >
            <span v-if="loading" class="absolute left-0 inset-y-0 flex items-center pl-3">
              <svg class="animate-spin h-5 w-5 text-indigo-300" fill="none" viewBox="0 0 24 24">
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"
                />
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                />
              </svg>
            </span>
            {{ loading ? '注册中...' : '立即注册' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, computed } from 'vue';
  import { useRouter } from 'vue-router';
  import { toast } from 'sonner';
  import { anixApi } from '@/api/anix';

  const router = useRouter();

  // 表单数据
  const form = ref({
    username: '',
    email: '',
    nickname: '',
    password: '',
    confirmPassword: '',
    agree: false,
  });

  // 表单错误
  const errors = ref<Record<string, string>>({});

  // 加载状态
  const loading = ref(false);

  // 密码强度计算
  const passwordStrength = computed(() => {
    const password = form.value.password;
    if (!password) return 0;

    let strength = 0;
    // 长度检查
    if (password.length >= 8) strength += 1;
    // 包含数字
    if (/\d/.test(password)) strength += 1;
    // 包含小写字母
    if (/[a-z]/.test(password)) strength += 1;
    // 包含大写字母
    if (/[A-Z]/.test(password)) strength += 1;
    // 包含特殊字符
    if (/[^\w\s]/.test(password)) strength += 1;

    return strength;
  });

  const passwordStrengthWidth = computed(() => {
    return `${(passwordStrength.value / 5) * 100}%`;
  });

  const passwordStrengthClass = computed(() => {
    const strength = passwordStrength.value;
    if (strength <= 1) return 'bg-red-500';
    if (strength <= 2) return 'bg-orange-500';
    if (strength <= 3) return 'bg-yellow-500';
    if (strength <= 4) return 'bg-blue-500';
    return 'bg-green-500';
  });

  const passwordStrengthTextClass = computed(() => {
    const strength = passwordStrength.value;
    if (strength <= 1) return 'text-red-600';
    if (strength <= 2) return 'text-orange-600';
    if (strength <= 3) return 'text-yellow-600';
    if (strength <= 4) return 'text-blue-600';
    return 'text-green-600';
  });

  const passwordStrengthText = computed(() => {
    const strength = passwordStrength.value;
    if (strength <= 1) return '弱';
    if (strength <= 2) return '较弱';
    if (strength <= 3) return '中等';
    if (strength <= 4) return '较强';
    return '强';
  });

  // 表单验证
  const validateForm = () => {
    errors.value = {};

    // 用户名验证
    if (!form.value.username.trim()) {
      errors.value.username = '请输入用户名';
    } else if (form.value.username.length < 3) {
      errors.value.username = '用户名至少3个字符';
    } else if (form.value.username.length > 20) {
      errors.value.username = '用户名不能超过20个字符';
    } else if (!/^[a-zA-Z0-9_]+$/.test(form.value.username)) {
      errors.value.username = '用户名只能包含字母、数字和下划线';
    }

    // 邮箱验证
    if (!form.value.email.trim()) {
      errors.value.email = '请输入邮箱地址';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
      errors.value.email = '请输入有效的邮箱地址';
    }

    // 昵称验证（可选）
    if (form.value.nickname && form.value.nickname.length > 50) {
      errors.value.nickname = '昵称不能超过50个字符';
    }

    // 密码验证
    if (!form.value.password) {
      errors.value.password = '请输入密码';
    } else if (form.value.password.length < 6) {
      errors.value.password = '密码至少6个字符';
    } else if (form.value.password.length > 50) {
      errors.value.password = '密码不能超过50个字符';
    }

    // 确认密码验证
    if (!form.value.confirmPassword) {
      errors.value.confirmPassword = '请确认密码';
    } else if (form.value.password !== form.value.confirmPassword) {
      errors.value.confirmPassword = '两次输入的密码不一致';
    }

    // 用户协议验证
    if (!form.value.agree) {
      errors.value.agree = '请阅读并同意用户协议和隐私政策';
    }

    return Object.keys(errors.value).length === 0;
  };

  // 处理注册
  const handleRegister = async () => {
    if (!validateForm()) {
      return;
    }

    loading.value = true;

    try {
      await anixApi.register({
        username: form.value.username,
        email: form.value.email,
        nickname: form.value.nickname || form.value.username,
        password: form.value.password,
      });

      toast.success('注册成功！请登录您的账户');
      router.push('/anix/login');
    } catch (error: any) {
      console.error('注册失败:', error);
      toast.error(error.message || '注册失败，请稍后重试');
    } finally {
      loading.value = false;
    }
  };
</script>
