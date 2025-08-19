<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- 页面标题 -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">个人中心</h1>
        <p class="mt-2 text-gray-600">管理您的个人信息和账户设置</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- 侧边栏导航 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow-sm border border-gray-200">
            <nav class="p-4 space-y-2">
              <button
                v-for="tab in tabs"
                :key="tab.key"
                @click="activeTab = tab.key"
                class="w-full text-left px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                :class="
                  activeTab === tab.key
                    ? 'bg-indigo-100 text-indigo-700'
                    : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
                "
              >
                <i :class="tab.icon" class="mr-3"></i>
                {{ tab.label }}
              </button>
            </nav>
          </div>
        </div>

        <!-- 主内容区 -->
        <div class="lg:col-span-2">
          <!-- 基本信息 -->
          <div
            v-if="activeTab === 'profile'"
            class="bg-white rounded-lg shadow-sm border border-gray-200"
          >
            <div class="p-6 border-b border-gray-200">
              <h2 class="text-lg font-medium text-gray-900">基本信息</h2>
              <p class="mt-1 text-sm text-gray-600">更新您的个人资料信息</p>
            </div>

            <form @submit.prevent="updateProfile" class="p-6 space-y-6">
              <!-- 头像上传 -->
              <div class="flex items-center space-x-6">
                <div class="shrink-0">
                  <img
                    class="h-20 w-20 object-cover rounded-full border-2 border-gray-200"
                    :src="profileForm.avatar || '/default-avatar.png'"
                    :alt="profileForm.nickname"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2"> 头像 </label>
                  <input
                    ref="avatarInput"
                    type="file"
                    accept="image/*"
                    @change="handleAvatarChange"
                    class="hidden"
                  />
                  <button
                    type="button"
                    @click="$refs.avatarInput.click()"
                    class="bg-white py-2 px-3 border border-gray-300 rounded-md shadow-sm text-sm leading-4 font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                  >
                    更换头像
                  </button>
                  <p class="mt-2 text-xs text-gray-500"> 支持 JPG、PNG 格式，文件大小不超过 2MB </p>
                </div>
              </div>

              <!-- 用户名 -->
              <div>
                <label for="username" class="block text-sm font-medium text-gray-700">
                  用户名
                </label>
                <input
                  id="username"
                  v-model="profileForm.username"
                  type="text"
                  disabled
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm bg-gray-50 text-gray-500 sm:text-sm"
                />
                <p class="mt-1 text-xs text-gray-500">用户名不可修改</p>
              </div>

              <!-- 昵称 -->
              <div>
                <label for="nickname" class="block text-sm font-medium text-gray-700"> 昵称 </label>
                <input
                  id="nickname"
                  v-model="profileForm.nickname"
                  type="text"
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  :class="{ 'border-red-500': errors.nickname }"
                />
                <p v-if="errors.nickname" class="mt-1 text-sm text-red-600">{{
                  errors.nickname
                }}</p>
              </div>

              <!-- 邮箱 -->
              <div>
                <label for="email" class="block text-sm font-medium text-gray-700">
                  邮箱地址
                </label>
                <input
                  id="email"
                  v-model="profileForm.email"
                  type="email"
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  :class="{ 'border-red-500': errors.email }"
                />
                <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
              </div>

              <!-- 个人简介 -->
              <div>
                <label for="bio" class="block text-sm font-medium text-gray-700"> 个人简介 </label>
                <textarea
                  id="bio"
                  v-model="profileForm.bio"
                  rows="3"
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="介绍一下自己吧..."
                ></textarea>
              </div>

              <!-- 保存按钮 -->
              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="profileLoading"
                  class="bg-indigo-600 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ profileLoading ? '保存中...' : '保存更改' }}
                </button>
              </div>
            </form>
          </div>

          <!-- 修改密码 -->
          <div
            v-if="activeTab === 'password'"
            class="bg-white rounded-lg shadow-sm border border-gray-200"
          >
            <div class="p-6 border-b border-gray-200">
              <h2 class="text-lg font-medium text-gray-900">修改密码</h2>
              <p class="mt-1 text-sm text-gray-600">确保您的账户安全</p>
            </div>

            <form @submit.prevent="changePassword" class="p-6 space-y-6">
              <!-- 当前密码 -->
              <div>
                <label for="currentPassword" class="block text-sm font-medium text-gray-700">
                  当前密码
                </label>
                <input
                  id="currentPassword"
                  v-model="passwordForm.currentPassword"
                  type="password"
                  required
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  :class="{ 'border-red-500': errors.currentPassword }"
                />
                <p v-if="errors.currentPassword" class="mt-1 text-sm text-red-600">{{
                  errors.currentPassword
                }}</p>
              </div>

              <!-- 新密码 -->
              <div>
                <label for="newPassword" class="block text-sm font-medium text-gray-700">
                  新密码
                </label>
                <input
                  id="newPassword"
                  v-model="passwordForm.newPassword"
                  type="password"
                  required
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  :class="{ 'border-red-500': errors.newPassword }"
                />
                <p v-if="errors.newPassword" class="mt-1 text-sm text-red-600">{{
                  errors.newPassword
                }}</p>
              </div>

              <!-- 确认新密码 -->
              <div>
                <label for="confirmPassword" class="block text-sm font-medium text-gray-700">
                  确认新密码
                </label>
                <input
                  id="confirmPassword"
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  required
                  class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  :class="{ 'border-red-500': errors.confirmPassword }"
                />
                <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{
                  errors.confirmPassword
                }}</p>
              </div>

              <!-- 保存按钮 -->
              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="passwordLoading"
                  class="bg-indigo-600 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ passwordLoading ? '修改中...' : '修改密码' }}
                </button>
              </div>
            </form>
          </div>

          <!-- 账户设置 -->
          <div
            v-if="activeTab === 'settings'"
            class="bg-white rounded-lg shadow-sm border border-gray-200"
          >
            <div class="p-6 border-b border-gray-200">
              <h2 class="text-lg font-medium text-gray-900">账户设置</h2>
              <p class="mt-1 text-sm text-gray-600">管理您的账户偏好设置</p>
            </div>

            <div class="p-6 space-y-6">
              <!-- 通知设置 -->
              <div>
                <h3 class="text-sm font-medium text-gray-900 mb-4">通知设置</h3>
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700">邮件通知</label>
                      <p class="text-xs text-gray-500">接收重要更新和通知邮件</p>
                    </div>
                    <input
                      v-model="settings.emailNotifications"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                  </div>
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700">推送通知</label>
                      <p class="text-xs text-gray-500">接收浏览器推送通知</p>
                    </div>
                    <input
                      v-model="settings.pushNotifications"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                  </div>
                </div>
              </div>

              <!-- 隐私设置 -->
              <div>
                <h3 class="text-sm font-medium text-gray-900 mb-4">隐私设置</h3>
                <div class="space-y-3">
                  <div class="flex items-center justify-between">
                    <div>
                      <label class="text-sm font-medium text-gray-700">公开个人资料</label>
                      <p class="text-xs text-gray-500">允许其他用户查看您的个人资料</p>
                    </div>
                    <input
                      v-model="settings.publicProfile"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    />
                  </div>
                </div>
              </div>

              <!-- 危险操作 -->
              <div class="border-t border-gray-200 pt-6">
                <h3 class="text-sm font-medium text-red-600 mb-4">危险操作</h3>
                <div class="bg-red-50 border border-red-200 rounded-md p-4">
                  <div class="flex">
                    <div class="flex-shrink-0">
                      <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                        <path
                          fill-rule="evenodd"
                          d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                          clip-rule="evenodd"
                        />
                      </svg>
                    </div>
                    <div class="ml-3">
                      <h3 class="text-sm font-medium text-red-800">删除账户</h3>
                      <div class="mt-2 text-sm text-red-700">
                        <p>一旦删除账户，所有数据将无法恢复。请谨慎操作。</p>
                      </div>
                      <div class="mt-4">
                        <button
                          type="button"
                          class="bg-red-600 py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                        >
                          删除账户
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref, reactive, onMounted } from 'vue';
  import { useMessage } from 'naive-ui';
  import { useAuthStore } from '@/stores/auth';
  import { anixApi } from '@/api/anix';

  const authStore = useAuthStore();
  const message = useMessage();

  // 当前激活的标签页
  const activeTab = ref('profile');

  // 标签页配置
  const tabs = [
    { key: 'profile', label: '基本信息', icon: 'fas fa-user' },
    { key: 'password', label: '修改密码', icon: 'fas fa-lock' },
    { key: 'settings', label: '账户设置', icon: 'fas fa-cog' },
  ];

  // 个人资料表单
  const profileForm = reactive({
    username: '',
    nickname: '',
    email: '',
    avatar: '',
    bio: '',
  });

  // 密码修改表单
  const passwordForm = reactive({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });

  // 设置
  const settings = reactive({
    emailNotifications: true,
    pushNotifications: false,
    publicProfile: true,
  });

  // 表单错误
  const errors = ref<Record<string, string>>({});

  // 加载状态
  const profileLoading = ref(false);
  const passwordLoading = ref(false);

  // 头像输入引用
  const avatarInput = ref<HTMLInputElement>();

  // 获取用户信息
  const fetchUserProfile = async () => {
    try {
      const user = await anixApi.getProfile();
      Object.assign(profileForm, user);
    } catch (error: any) {
      console.error('获取用户信息失败:', error);
      message.error('获取用户信息失败');
    }
  };

  // 处理头像上传
  const handleAvatarChange = async (event: Event) => {
    const file = (event.target as HTMLInputElement).files?.[0];
    if (!file) return;

    // 验证文件类型
    if (!file.type.startsWith('image/')) {
      message.error('请选择图片文件');
      return;
    }

    // 验证文件大小（2MB）
    if (file.size > 2 * 1024 * 1024) {
      message.error('图片大小不能超过 2MB');
      return;
    }

    try {
      const formData = new FormData();
      formData.append('avatar', file);

      const response = await anixApi.uploadAvatar(formData);
      profileForm.avatar = response.url;
      message.success('头像上传成功');
    } catch (error: any) {
      console.error('头像上传失败:', error);
      message.error('头像上传失败');
    }
  };

  // 验证个人资料表单
  const validateProfileForm = () => {
    errors.value = {};

    if (!profileForm.nickname.trim()) {
      errors.value.nickname = '请输入昵称';
    } else if (profileForm.nickname.length > 50) {
      errors.value.nickname = '昵称不能超过50个字符';
    }

    if (!profileForm.email.trim()) {
      errors.value.email = '请输入邮箱地址';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(profileForm.email)) {
      errors.value.email = '请输入有效的邮箱地址';
    }

    return Object.keys(errors.value).length === 0;
  };

  // 验证密码表单
  const validatePasswordForm = () => {
    errors.value = {};

    if (!passwordForm.currentPassword) {
      errors.value.currentPassword = '请输入当前密码';
    }

    if (!passwordForm.newPassword) {
      errors.value.newPassword = '请输入新密码';
    } else if (passwordForm.newPassword.length < 6) {
      errors.value.newPassword = '新密码至少6个字符';
    }

    if (!passwordForm.confirmPassword) {
      errors.value.confirmPassword = '请确认新密码';
    } else if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      errors.value.confirmPassword = '两次输入的密码不一致';
    }

    return Object.keys(errors.value).length === 0;
  };

  // 更新个人资料
  const updateProfile = async () => {
    if (!validateProfileForm()) {
      return;
    }

    profileLoading.value = true;

    try {
      await anixApi.updateProfile({
        nickname: profileForm.nickname,
        email: profileForm.email,
        avatar: profileForm.avatar,
        bio: profileForm.bio,
      });

      // 更新本地用户信息
      authStore.updateUser({
        nickname: profileForm.nickname,
        email: profileForm.email,
        avatar: profileForm.avatar,
        bio: profileForm.bio,
      });

      message.success('个人资料更新成功');
    } catch (error: any) {
      console.error('更新个人资料失败:', error);
      message.error(error.message || '更新失败，请稍后重试');
    } finally {
      profileLoading.value = false;
    }
  };

  // 修改密码
  const changePassword = async () => {
    if (!validatePasswordForm()) {
      return;
    }

    passwordLoading.value = true;

    try {
      await anixApi.changePassword({
        currentPassword: passwordForm.currentPassword,
        newPassword: passwordForm.newPassword,
      });

      // 清空表单
      passwordForm.currentPassword = '';
      passwordForm.newPassword = '';
      passwordForm.confirmPassword = '';

      message.success('密码修改成功');
    } catch (error: any) {
      console.error('修改密码失败:', error);
      message.error(error.message || '修改密码失败，请稍后重试');
    } finally {
      passwordLoading.value = false;
    }
  };

  // 组件挂载时获取用户信息
  onMounted(() => {
    fetchUserProfile();
  });
</script>
