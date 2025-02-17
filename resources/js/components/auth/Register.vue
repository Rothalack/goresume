<script>
	import { ref, watch, computed, onMounted } from 'vue'
	export default {
		setup() {
			const userEmail = ref('')
			const userName = ref('')
			const password = ref('')
			const confirm = ref('')
			const isLoading = ref(false)
			const errors = ref({})

			const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    		const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{8,}$/;

			const validations = computed(() => ({
				email: {
					isValid: emailRegex.test(userEmail.value),
					message: 'Please enter a valid email address'
				},
				password: {
					isValid: passwordRegex.test(password.value),
					message: 'Password must be at least 8 characters and contain both letters and numbers'
				},
				passwordMatch: {
					isValid: password.value === confirm.value,
					message: 'Passwords do not match'
				},
				userName: {
					isValid: userName.value.length >= 2,
					message: 'Name must be at least 2 characters long'
				}
			}));

			const isFormValid = computed(() =>
				Object.values(validations.value).every(v => v.isValid) &&
				userEmail.value && userName.value && password.value && confirm.value
			);

			const register = async () => {
				if (!isFormValid.value) return;

				isLoading.value = true;
				errors.value = {};

				try {
					const response = await fetch('/auth/register', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						user_email: userEmail.value,
						password: password.value,
						user_name: userName.value
					})
					});

					const data = await response.json();

					if (!response.ok) {
						throw new Error(data.error || 'Registration failed');
					}

					window.location.href = '/';
				} catch (error) {
					errors.value.submit = error.message;
				} finally {
					isLoading.value = false;
				}
			};

			return {
				userEmail,
				userName,
				password,
				confirm,
				isLoading,
				errors,
				validations,
				isFormValid,
				register
			};
		},
	}
</script>
<template>
	<div class="max-w-sm mx-auto mt-5">
		<h2 class="text-2xl font-bold text-white mb-6">Create Account</h2>
		<div v-if="errors.submit" class="mb-4 p-4 bg-red-500 text-white rounded-lg">
			{{ errors.submit }}
		</div>

		<form @submit.prevent="register" class="space-y-4">
			<div id="useremail-input">
				<label for="useremail" class="block mb-2 text-sm font-medium text-white">
					Email Address
				</label>
				<input v-model="userEmail" type="email" id="useremail" placeholder="you@example.com" required
					:class="[
						'bg-gray-50 border text-gray-900 text-sm rounded-lg block w-full p-2.5',
						userEmail && !validations.email.isValid ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'
					]" />
				<p v-if="userEmail && !validations.email.isValid" class="mt-1 text-sm text-red-500">
					{{ validations.email.message }}
				</p>
			</div>
			<div id="username-input">
				<label for="username" class="block mb-2 text-sm font-medium text-white">
					User Name
				</label>
				<input v-model="userName" type="text" id="username" placeholder="John Doe" required
					:class="[
						'bg-gray-50 border text-gray-900 text-sm rounded-lg block w-full p-2.5',
						userName && !validations.userName.isValid ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'
					]"/>
				<p v-if="userName && !validations.userName.isValid" class="mt-1 text-sm text-red-500">
					{{ validations.userName.message }}
				</p>
			</div>
			<div id="password-input">
				<label for="password" class="block mb-2 text-sm font-medium text-white">
					Password
				</label>
				<input v-model="password" type="password" id="password" placeholder="*********" required
					:class="[
						'bg-gray-50 border text-gray-900 text-sm rounded-lg block w-full p-2.5',
						password && !validations.password.isValid ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'
					]" />
				<p v-if="password && !validations.password.isValid" class="mt-1 text-sm text-red-400">
					{{ validations.password.message }}
				</p>
			</div>
			<div id="confirm-input">
				<label for="confirm" class="block mb-2 text-sm font-medium text-white">
					Confirm Password
				</label>
				<input v-model="confirm" type="password" id="confirm" placeholder="••••••••" required
					:class="[
						'bg-gray-50 border text-gray-900 text-sm rounded-lg block w-full p-2.5',
						confirm && !validations.passwordMatch.isValid ? 'border-red-500 focus:ring-red-500 focus:border-red-500' : 'border-gray-300 focus:ring-blue-500 focus:border-blue-500'
					]"
				/>
				<p v-if="confirm && !validations.passwordMatch.isValid" class="mt-1 text-sm text-red-500">
					{{ validations.passwordMatch.message }}
				</p>
			</div>

			<button type="submit" :disabled="!isFormValid || isLoading" :class="[
					'w-full px-4 py-2 text-white font-semibold rounded-lg shadow-lg transition',
					isFormValid && !isLoading ? 'bg-blue-600 hover:bg-blue-700' : 'bg-gray-400 cursor-not-allowed'
			]">
				<span v-if="isLoading" class="inline-flex items-center">
					<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
						<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
						<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
					</svg>
					Creating account...
				</span>
				<span v-else>Create Account</span>
			</button>

			<p class="text-sm text-gray-400 mt-4 text-center">
				Already have an account?
				<a href="/login" class="text-blue-500 hover:text-blue-600">Sign in</a>
			</p>
		</form>
	</div>
</template>