import { autoRegisterComponents } from '/resources/js/utils/registerComponents.js';
import GuildRankings from '/resources/js/components/GuildRankings.vue';
import Login from '/resources/js/components/auth/Login.vue';
import Register from '/resources/js/components/auth/Register.vue';

autoRegisterComponents({
	'guild-rankings': GuildRankings,
	'login-form': Login,
	'register-form': Register,
});
