const path = require('path');
const { VueLoaderPlugin } = require('vue-loader');

module.exports = {
	entry: './resources/js/app.js',
	output: {
		filename: 'bundle.js',
		path: path.resolve(__dirname, 'static/js'),
	},
	resolve: {
		alias: {
			vue: 'vue/dist/vue.esm-browser.js',
		},
		extensions: ['.js', '.vue', '.json'],
	},
	module: {
		rules: [
			{
				test: /\.vue$/,
				use: 'vue-loader',
			},
		],
	},
	plugins: [new VueLoaderPlugin()],
};
