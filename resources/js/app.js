import { autoRegisterComponents } from '/resources/js/utils/registerComponents.js';
import ExampleTest from '/resources/js/components/ExampleTest.vue';
import GuildRankings from '/resources/js/components/GuildRankings.vue';
import CharacterRankings from '/resources/js/components/CharacterRankings.vue';
import VueSelect from 'vue3-select-component';

autoRegisterComponents({
  'example-test': ExampleTest,
  'guild-rankings': GuildRankings,
  'character-rankings': CharacterRankings,
  'vue-select': VueSelect,
});