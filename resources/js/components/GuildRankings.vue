<script>
import { ref, watch, computed, onMounted } from 'vue';

export default {
	setup() {
		const loading = ref(false)
		const logsData = ref(null)
		const rankingData = ref(null)
		const charData =ref([])
		const guildName = ref(null)
		const guildId = ref(null)
		const guildFaction = ref(null)
		const selectedGame = ref(null)
		const selectedRegion = ref(null)
		const selectedServer = ref(null)
		const selectedExpansion = ref(null)
		const selectedZone = ref(null)
		const selectedDifficulty = ref(null)
		const selectedSize = ref(null)
		const errorMessage = ref(null)

		const fetchData = async () => {
			loading.value = true
			try {
				const response = await fetch('/api/logs-data')
				const data = await response.json()
				logsData.value = data.data
			} catch (error) {
				console.error('Error fetching logs data:', error)
				errorMessage.value = 'Error fetching logs data:' + error
			} finally {
				loading.value = false;
			}
		}

		const fetchRanking = async () => {
			loading.value = true
			try {
				const params = new URLSearchParams({
					guild: guildName.value,
					api_url: selectedGame.value.api_url,
					server: selectedServer.value.server_name,
					region: selectedRegion.value.slug,
					zone: selectedZone.value.zone_id,
					difficulty: selectedDifficulty.value?.difficulty_id ?? 0,
					size: selectedSize.value?.size ?? 0,
				})
				const response = await fetch(`/api/ranking-data?${params.toString()}`)
				const data = await response.json()

				if (!response.ok) {
					throw new Error(data.error || 'Unknown error occurred');
				}

				rankingData.value = data.data
				guildId.value = data.guildId
				guildFaction.value = data.guildFaction
			} catch (error) {
				console.error('Error fetching ranking data: ', error)
				errorMessage.value = 'Error fetching ranking data: ' + error
			} finally {
				loading.value = false;
			}
			fetchChars()
		}

		const fetchChars = async () => {
			loading.value = true
			try {
				const params = new URLSearchParams({
					guild: guildId.value,
					zone: selectedZone.value.zone_id
				})
				const response = await fetch(`/api/char-data?${params.toString()}`)
				const data = await response.json()
				charData.value = data.data
			} catch (error) {
				console.error('Error fetching character data: ', error)
				errorMessage.value = 'Error fetching character data: ' + error
			} finally {
				loading.value = false;
			}
		}

		const classMap = computed(() => {
			if (logsData.value && Array.isArray(logsData.value[0].classes)) {
				return Object.fromEntries(logsData.value[0].classes.map(cls => [cls.class_id, cls.class_name.replace(' ', '')]));
			}
			return {};
		});

		const formatNumberSuffix = (number) => {
			const j = number % 10
			const k = number % 100
			if (j === 1 && k !== 11) {
				return `${number}st`
			}
			if (j === 2 && k !== 12) {
				return `${number}nd`
			}
			if (j === 3 && k !== 13) {
				return `${number}rd`
			}
			return `${number}th`
		}

		const resetForm = () => {
			rankingData.value = null
			charData.value = null
		}

		onMounted(() => {
			fetchData()
		});

		watch(
			[guildName, selectedGame, selectedRegion, selectedServer, selectedExpansion, selectedZone, selectedDifficulty, selectedSize],
			() => {
				errorMessage.value = null;
			}
		);

		return {
			loading,
			logsData,
			rankingData,
			charData,
			guildName,
			guildFaction,
			selectedGame,
			selectedRegion,
			selectedServer,
			selectedExpansion,
			selectedZone,
			selectedDifficulty,
			selectedSize,
			fetchData,
			fetchRanking,
			fetchChars,
			formatNumberSuffix,
			classMap,
			errorMessage,
			resetForm,
		}
	},
}
</script>

<template>
	<div v-if="errorMessage" class="bg-red-500 text-white p-4 rounded text-center w-md flex justify-center">
		{{ errorMessage }}
	</div>
	<div id="guild_display_wrapper" v-if="rankingData">
		<button @click="resetForm" class="reset-btn absolute p-4 bg-gray-800 hover:bg-gray-700 focus:ring-1 focus:ring-gray-300 bg-opacity-50 top-1 right-1 lg:top-2 lg:right-2 rounded text-white hover:bg-opacity-75 z-10">
			<i class="fas fa-times"></i>
		</button>
		<div id="guild_display" class="lg:bg-contain bg-fixed bg-center bg-no-repeat bg-cover" :style="{ backgroundImage: `url('./static/images/raid_backgrounds/zone_${selectedZone.zone_id}.jpg')` }"></div>
		<div id="guild_content">
			<div class="flex justify-center transparent p-4">
				<div class="gap-4 max-w-4xl w-full p-4 bg-gray-800 shadow rounded">
					<h2 class="text-4xl mt-5 mb-5 text-center">{{ selectedZone.zone_name }}</h2>
					<h2 class="text-3xl mt-5 mb-5 text-center"><span :class="guildFaction"><strong>{{ guildName }}</strong></span> Guild Ranking</h2>
				</div>
			</div>
			<div class="flex justify-center transparent p-4">
				<div class="grid grid-cols-1 sm:grid-cols-3 gap-4 max-w-4xl w-full">
					<div class="p-4 bg-gray-800 shadow rounded text-center">
						<h2 class="text-3xl mt-4">Progression</h2>
						<p class="mt-4"><span class="text-xl">World</span><span :class="rankingData.progress.worldRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.progress.worldRank.number) }}</span></p>
						<p class="mt-4"><span class="text-xl">Region</span><span :class="rankingData.progress.regionRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.progress.regionRank.number) }}</span></p>
						<p class="mt-4 mb-4"><span class="text-xl">Server</span><span :class="rankingData.progress.serverRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.progress.serverRank.number) }}</span></p>
					</div>
					<div class="p-4 bg-gray-800 shadow rounded text-center">
						<h2 class="text-3xl mt-4">Speed Runs</h2>
						<p class="mt-4"><span class="text-xl">World</span><span :class="rankingData.completeRaidSpeed.worldRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.completeRaidSpeed.worldRank.number) }}</span></p>
						<p class="mt-4"><span class="text-xl">Region</span><span :class="rankingData.completeRaidSpeed.regionRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.completeRaidSpeed.regionRank.number) }}</span></p>
						<p class="mt-4 mb-4"><span class="text-xl">Server</span><span :class="rankingData.completeRaidSpeed.serverRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.completeRaidSpeed.serverRank.number) }}</span></p>
					</div>
					<div class="p-4 bg-gray-800 shadow rounded text-center">
						<h2 class="text-3xl mt-4">Fight Speeds</h2>
						<p class="mt-4"><span class="text-xl">World</span><span :class="rankingData.speed.worldRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.speed.worldRank.number) }}</span></p>
						<p class="mt-4"><span class="text-xl">Region</span><span :class="rankingData.speed.regionRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.speed.regionRank.number) }}</span></p>
						<p class="mt-4 mb-4"><span class="text-xl">Server</span><span :class="rankingData.speed.serverRank.color" class="text-2xl ml-2">{{ formatNumberSuffix(rankingData.speed.serverRank.number) }}</span></p>
					</div>
				</div>
			</div>
			<h2 class="text-4xl mt-5 mb-5 text-center">Roster</h2>
			<div class="flex justify-center transparent">
				<div class="max-w-4xl w-full">
					<div v-for="(charDatum, charIndex) in charData" class="max-w-4xl w-full">
						<div :class="['p-4', charIndex % 2 === 0 ? 'bg-gray-800' : 'bg-gray-700']">
							{{ charDatum.level }} <span :class="classMap[charDatum.classId]">{{ charDatum.name }}</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
	<div class="relative">
		<div v-if="loading" class="loading-overlay">
			<div class="loading-spinner"></div>
		</div>
		<div class="max-w-sm mx-auto mt-5" v-if="!rankingData">
			<h3 class="text-2xl mb-4">Search for a WoW Guild Ranking</h3>
			<div class="mb-5" id="guild-input">
				<label for="guild" class="block mb-2 text-sm font-medium text-white dark:text-white">
					Guild Name
				</label>
				<input
					v-model="guildName"
					type="text"
					id="guild"
					class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
					placeholder="Enter guild name"
					required
				/>
			</div>

			<div class="mb-5" id="game-input">
				<label for="game" class="block mb-2 text-sm font-medium text-white dark:text-white">Game</label>
				<select id="game" v-model="selectedGame" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in logsData" :key="option.game_id" :value="option">
						{{ option.game_name }}
					</option>
				</select>
			</div>

			<div class="mb-5" id="region-input" v-if="selectedGame">
				<label for="region" class="block mb-2 text-sm font-medium text-white dark:text-white">Region</label>
				<select id="region" v-model="selectedRegion" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in selectedGame.regions" :key="option.region_id" :value="option">
						{{ option.compact_name }}
					</option>
				</select>
			</div>

			<div class="mb-5" id="server-input" v-if="selectedRegion">
				<label for="server" class="block mb-2 text-sm font-medium text-white dark:text-white">Server</label>
				<select id="server" v-model="selectedServer" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in selectedRegion.servers" :key="option.server_id" :value="option">
						{{ option.normalized_name }}
					</option>
				</select>
			</div>

			<div class="mb-5" id="expansion-input" v-if="selectedServer">
				<label for="expansion" class="block mb-2 text-sm font-medium text-white dark:text-white">
					Expansion
				</label>
				<select id="expansion" v-model="selectedExpansion" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in selectedGame.expansions" :key="option.expansion_id" :value="option">
						{{ option.expansion_name }}
					</option>
				</select>
			</div>

			<div class="mb-5" id="zone-input" v-if="selectedExpansion">
				<label for="zone" class="block mb-2 text-sm font-medium text-white dark:text-white">Zone</label>
				<select id="expansion" v-model="selectedZone" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in selectedExpansion.zones" :key="option.zone_id" :value="option">
						{{ option.zone_name }}
					</option>
				</select>
			</div>

			<div class="mb-5" id="difficulty-input" v-if="selectedZone && selectedZone.difficulty && selectedZone.difficulty.length > 1">
				<label for="difficulty" class="block mb-2 text-sm font-medium text-white dark:text-white">
					Difficulty
				</label>
				<select id="difficulty" v-model="selectedDifficulty" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in selectedZone.difficulty" :key="option.difficulty_id" :value="option">
						{{ option.difficulty_name }}
					</option>
				</select>
			</div>

			<div class="mb-5" id="size-input" v-if="selectedDifficulty && selectedDifficulty.sizes && selectedDifficulty.sizes.length > 1">
				<label for="size" class="block mb-2 text-sm font-medium text-white dark:text-white">Size</label>
				<select id="difficulty" v-model="selectedSize" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
					<option v-for="option in selectedDifficulty.sizes" :key="option.size" :value="option.size">
						{{ option.size }}
					</option>
				</select>
			</div>

			<button v-if="guildName && selectedZone"
				@click="fetchRanking"
				id="search-rankings"
				class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800"
				:disabled="!guildName || !selectedZone"
			>
				Search Rankings
			</button>
		</div>
	</div>
</template>
