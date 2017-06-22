(() => {"use strict";




const REGEX_SENTENSE = /[^\s].+?[.!?\s]+(?=\s)/g;
const REGEX_ENDLINE = /\s*\n\s+/g;
const REGEX_INVISIBLE_CHARS = /[\t\n\s]+/g;
const REGEX_QUOTES = /«|»/g;

const STATS = {
    time: 0,
    speed: 0,
    av: 0,
    acc: 100,
    max: 0,
    total: 0,
    nErrors: 0
}




window.app = new Vue({
    el: '#vue_container',

    data: {

        stats: Object.assign({}, STATS),
        localeFace: {},

        settings: {
            stamina: false,

            timer: {
                enabled: false,
                time: 60,
                func: 0,
                alert: false
            },

            randomText: true,

            locale: {
                faceLang: localStorage.face || 'en',
                textsLang: localStorage.texts || 'en',
                list: ['en', 'ru']
            }
        },

        textSource: '',
        text: '',
        distance: 17,
        pos: 0,
        errorFull: '',
        taps: [],

        stage: 0,
        startTime: 0,
        updateInterval: 0,

        settingsMode: false,
        editMode: false
    },

    methods: {
        onKeyDown(event) {
            if ((this.stage === 2) ||
                (event.key.length > 1 && event.keyCode !== 8) ||
                (!this.text) ||
                (this.editMode) ||
                (this.settingsMode)
            ) return;

            if (this.stage === 0) this.start();

            event.preventDefault();


            if (event.keyCode === 8) {
                if (!this.settings.stamina) return this.erease();
                return;
            }

            if (event.key !== this.text[this.pos] || this.errorFull) {
                this.stats.nErrors++;
                if (!this.settings.stamina) this.errorFull += event.key;
            } else {
                this.taps.push(performance.now());
                this.pos++;
                if (this.pos === this.text.length) this.stop();
            }

            this.stats.total++;
            this.updateStats();
        },

        erease() {
            if (!this.errorFull) return (this.pos) ? this.pos-- : null;
            this.errorFull = this.errorFull.slice(0, -1);
        },

        start() {
            this.startTime = performance.now();
            this.updateInterval = setInterval(this.updateStats, 1000);
            this.stage = 1;
            
            let { timer } = this.settings;
            if (timer.enabled) {
                timer.func = setTimeout(this.timerAlert, timer.time * 1000);
            }
        },

        updateStats() {
            let now = performance.now();
            let { stats, startTime, updateTaps } = this;
            
            updateTaps(now);

            let valid = stats.total - stats.nErrors;
            let av = Math.floor(valid / stats.time * 60);
            let acc = 100 / (stats.total / (stats.total - stats.nErrors));

            stats.speed = this.taps.length;
            stats.time = Math.round((now - startTime) / 1000);
            stats.av = (av < 2000) ? av : 2000;
            stats.acc = isFinite(acc) ? Math.floor(acc) : 0;
            if (stats.max < stats.speed) stats.max = stats.speed;
        },

        updateTaps(now) {
            let ind = this.taps.findIndex(t =>  t + 60000 > now);
            if (ind) this.taps = (ind > -1) ? this.taps.slice(ind) : [];
        },

        stop() {
            if (this.stage === 1) {
                clearInterval(this.updateInterval);
                clearTimeout(this.settings.timer.func);
                this.stage = 2;
                this.updateStats();
            }
        },

        reset() {
            if (!this.stage) return;
            clearInterval(this.updateInterval);
            clearTimeout(this.settings.timer.func);
            this.stage = this.pos = 0;
            this.errorFull = '';
            this.taps = [];
            Object.assign(this.stats, STATS);
        },

        timerAlert() {
            this.stop();
            this.settings.timer.alert = true;
            setTimeout(() => { this.settings.timer.alert = false; }, 1000);
        },

        randomize(text) {
            return text
                .match(REGEX_SENTENSE)
                .sort(() => Math.random() - .5)
                .join(' ')
                .replace(REGEX_ENDLINE, '\n');
        },

        async newText() {
            let locale = this.settings.locale.textsLang;
            let res = await fetch(`/random_text/?locale=${locale}`);
            let text = await res.text();
            let isRand = this.settings.randomText;
            this.textSource = isRand ? this.randomize(text) : text;
        },

        async newLocale() {
            let locale = this.settings.locale.faceLang;
            let res = await fetch(`/locales/${locale}.json`);
            this.localeFace = await res.json();
        }
    },

    computed: {
        pass() {
            let { text, pos, distance, errorFull } = this;
            let from = pos - distance + errorFull.length;
            let chars = (from < pos) ? text.substring(from, pos) : '';
            let nSpaces = distance - chars.length - errorFull.length;
            return ' '.repeat((nSpaces > 0) ? nSpaces : 0) + chars;
        },

        error() {
            let from = this.errorFull.length - this.distance;
            return this.errorFull.substring(from);
        }
    },

    watch: {
        textSource(t) {
            this.reset();
            this.text = t
                .replace(REGEX_INVISIBLE_CHARS, ' ')
                .replace(REGEX_QUOTES, '"')
                .replace('—', '-');
        },

        'settings.locale.faceLang': function(val) {
            localStorage.face = val;
            this.newLocale();
        },

        'settings.locale.textsLang': function(val) {
            localStorage.texts = val;
            this.newText();
        },

        'settings.randomText': function(val) {
            if (val) this.textSource = this.randomize(this.textSource);
        }
    },

    mounted() {
        this.$watch('settings.stamina', this.reset);
        this.$watch('settings.timer.enabled', this.reset);
        this.$watch('settings.randomText', this.reset);

        this.newLocale();
        this.newText();
        window.addEventListener('keydown', this.onKeyDown);
    }
});


})();
