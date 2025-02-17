<script>
	import { ref, watch, computed, onMounted } from 'vue'
	export default {
		setup() {
			const userEmail = ref('')
			const password = ref('')
			const isLoading = ref(false)
			const errors = ref({})

			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
		},
	}
</script>

<template>
	<div class="max-w-sm mx-auto mt-5">
		<div class="mb-5" id="user-input">
			<label for="user" class="block mb-2 text-sm font-medium text-white dark:text-white">User Email</label>
			<input v-model="userEmail" type="text" id="user" placeholder="Enter User Email" required
				class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"/>
		</div>
		<div class="mb-5" id="password-input">
			<label for="password" class="block mb-2 text-sm font-medium text-white dark:text-white">Password</label>
			<input v-model="password" type="text" id="password" placeholder="Enter Password" required
				class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"/>
		</div>
		<button @click="login" id="login" :disabled="!userEmail && !password"
			class="inline-flex items-center gap-2 px-4 py-2 bg-slate-700 text-white font-semibold rounded-lg shadow-lg hover:bg-slate-600 transition">
			Login
		</button>
	</div>
</template>