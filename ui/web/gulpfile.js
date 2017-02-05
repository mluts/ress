var gulp    = require('gulp'),
    webpack = require('webpack-stream'),
    sass    = require('gulp-sass'),
    exec    = require('child_process').exec,
    webpackCfg = require('./webpack.config.js');

function logError(err) {
  console.error(err);
}

gulp.task('webpack', function() {
  return gulp.src('src/index.js')
    .pipe(webpack(webpackCfg))
    .pipe(gulp.dest('static/assets/js'));
});

gulp.task('webpack:watch', function() {
  webpackCfg.watch = true;

  return gulp.src('src/index.js')
    .pipe(webpack(webpackCfg))
    .on('error', logError)
    .pipe(gulp.dest('static/assets/js'));
});

gulp.task('sass', function() {
  return gulp.src('src/scss/app.scss')
    .pipe(sass({includePaths: ['./src/scss']})
    .on('error', sass.logError))
    .pipe(gulp.dest('./static/assets/css'));
});

gulp.task('sass:watch', function() {
  gulp.watch('./src/scss/*.scss', ['sass']);
});

gulp.task('default', ['webpack', 'sass']);
gulp.task('watch', ['webpack:watch', 'sass:watch']);
