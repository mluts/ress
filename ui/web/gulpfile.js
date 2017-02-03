var gulp    = require('gulp'),
    webpack = require('webpack-stream'),
    sass    = require('gulp-sass'),
    exec    = require('child_process').exec;

gulp.task('webpack', function() {
  return gulp.src('src/index.js')
    .pipe(webpack({
      output: {
        filename: "app.js"
      }
    }))
    .pipe(gulp.dest('static/assets/js'));
});

gulp.task('webpack:watch', function() {
  return gulp.src('src/index.js')
    .pipe(webpack({
      watch: true,
      output: {
        filename: "app.js"
      }
    }))
    .pipe(gulp.dest('static/assets/js'));
});

gulp.task('sass', function() {
  return gulp.src('src/scss/app.scss')
    .pipe(sass({includePaths: ['./src/scss']}).on('error', sass.logError))
    .pipe(gulp.dest('./static/assets/css'));
});

gulp.task('sass:watch', function() {
  gulp.watch('./src/scss/*.scss', ['sass']);
});

gulp.task('ress', function(cb) {
  exec('ress', function(err) {
    if(err) return cb(err);
    cb();
  });
});

gulp.task('default', ['webpack', 'sass']);
gulp.task('watch', ['webpack:watch', 'sass:watch', 'ress']);
