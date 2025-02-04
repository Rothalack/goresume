<script>
import { ref, onMounted } from 'vue';

export default {
	setup() {
		const logsData = ref(null)
		const rankingData = ref(null)
		const charData =ref([])
		const guildName = ref(null)
		const guildId = ref(null)
		const selectedGame = ref(null)
		const selectedRegion = ref(null)
		const selectedServer = ref(null)
		const selectedExpansion = ref(null)
		const selectedZone = ref(null)
		const selectedDifficulty = ref(null)
		const selectedSize = ref(null)
		// const addingChar = ref(false)
		// const charName = ref('')

		const fetchData = async () => {
			try {
				const response = await fetch('/api/logs-data')
				const data = await response.json()
				logsData.value = data.data
			} catch (error) {
				console.error('Error fetching logs data:', error)
			}
		}

		const fetchRanking = async () => {
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
				rankingData.value = data.data
				guildId.value = data.guildId
			} catch (error) {
				console.error('Error fetching ranking data: ', error)
			}
			fetchChars()
		}

		const fetchChars = async () => {
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
			}
		}

		// const addChar = () => {
		// 	if (charName.value.trim() !== '') {
		// 		fetchChars()
		// 		charData.value.push(charName.value.trim())
		// 		charName.value = ''
		// 		addingChar.value = false
		// 	}
		// }

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

		onMounted(() => {
			fetchData()
		});

		return {
			logsData,
			rankingData,
			charData,
			guildName,
			selectedGame,
			selectedRegion,
			selectedServer,
			selectedExpansion,
			selectedZone,
			selectedDifficulty,
			selectedSize,
			// addingChar,
			// charName,
			fetchData,
			fetchRanking,
			fetchChars,
			formatNumberSuffix,
			// addChar,
		}
	},
}
</script>

<template>
	<div id="guild_display" class="shadow rounded-xl" :style="{ backgroundImage: `url('./static/images/raid_backgrounds/zone_${selectedZone.zone_id}.jpg')` }" v-if="rankingData">
		<h2 class="text-4xl mt-5 mb-5 text-center">Guild Ranking</h2>
		<div class="flex justify-center transparent">
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-4 max-w-4xl w-full p-4">
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
		<h2 class="text-4xl mt-5 mb-5 text-center">Character Rankings</h2>
		<div class="flex justify-center transparent">
			<div class="gap-4 max-w-4xl w-full p-4">
				<div v-for="(charDatum, charIndex) in charData" class="grid grid-cols-1 max-w-4xl w-full">
					<div :class="['p-2 text-center', charIndex % 2 === 0 ? 'bg-gray-800' : 'bg-gray-700']">
						{{ charIndex }} {{ charDatum }}
					</div>
				</div>
				<!-- <div class="grid grid-cols-1 max-w-sm mx-auto mt-5" v-if="rankingData">
					<button
						v-if="!addingChar"
						@click="addingChar = true"
						class="text-white bg-gray-600 hover:bg-gray-700 focus:ring-4 focus:outline-none focus:ring-blue-300 rounded-lg w-full px-5 py-2.5 text-center dark:bg-gray-600 dark:hover:bg-gray-700 dark:focus:ring-blue-800">
						<i class="fa-solid fa-plus fa-lg"></i>
					</button>
					<div v-if="addingChar">
						<label for="charName" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
							Character Name
						</label>
						<input
							v-model="charName"
							type="text"
							id="charName"
							class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
							placeholder="Enter character name"
							required
						/>
						<button
							@click="addChar"
							class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm w-full sm:w-auto px-5 py-2.5 text-center dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">
							Add Character
						</button>
					</div>
				</div> -->
			</div>
		</div>
	</div>
	<div class="max-w-sm mx-auto mt-5" v-if="!rankingData">
		<div class="mb-5" id="guild-input">
			<label for="guild" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
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
			<label for="game" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Game</label>
			<select id="game" v-model="selectedGame" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
				<option v-for="option in logsData" :key="option.game_id" :value="option">
					{{ option.game_name }}
				</option>
			</select>
		</div>

		<div class="mb-5" id="region-input" v-if="selectedGame">
			<label for="region" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Region</label>
			<select id="region" v-model="selectedRegion" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
				<option v-for="option in selectedGame.regions" :key="option.region_id" :value="option">
					{{ option.compact_name }}
				</option>
			</select>
		</div>

		<div class="mb-5" id="server-input" v-if="selectedRegion">
			<label for="server" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Server</label>
			<select id="server" v-model="selectedServer" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
				<option v-for="option in selectedRegion.servers" :key="option.server_id" :value="option">
					{{ option.normalized_name }}
				</option>
			</select>
		</div>

		<div class="mb-5" id="expansion-input" v-if="selectedServer">
			<label for="expansion" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
				Expansion
			</label>
			<select id="expansion" v-model="selectedExpansion" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
				<option v-for="option in selectedGame.expansions" :key="option.expansion_id" :value="option">
					{{ option.expansion_name }}
				</option>
			</select>
		</div>

		<div class="mb-5" id="zone-input" v-if="selectedExpansion">
			<label for="zone" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Zone</label>
			<select id="expansion" v-model="selectedZone" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
				<option v-for="option in selectedExpansion.zones" :key="option.zone_id" :value="option">
					{{ option.zone_name }}
				</option>
			</select>
		</div>

		<div class="mb-5" id="difficulty-input" v-if="selectedZone && selectedZone.difficulty">
			<label for="difficulty" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">
				Difficulty
			</label>
			<select id="difficulty" v-model="selectedDifficulty" class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
				<option v-for="option in selectedZone.difficulty" :key="option.difficulty_id" :value="option">
					{{ option.difficulty_name }}
				</option>
			</select>
		</div>

		<div class="mb-5" id="size-input" v-if="selectedDifficulty && selectedDifficulty.sizes">
			<label for="size" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">Size</label>
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
</template>
