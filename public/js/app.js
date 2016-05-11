/**
 * Created by tugceakin on 3/13/16.
 */
//var bouncerApp = angular.module('bouncerApp', ['chart.js', 'ngRoute'])

var bouncerApp = angular.module('bouncerApp', ['chart.js', 'ngRoute', 'ngAnimate', 'uiSwitch', 'ui.bootstrap'])

    .config(function ($routeProvider) {
    $routeProvider.
    when('/', {
        templateUrl: 'benchmarking.html',
        controller: 'BenchmarkingController'
    }).
    when('/benchmarking', {
        templateUrl: 'benchmarking.html',
        controller: 'BenchmarkingController'
    }).
    when('/configurations', {
        templateUrl: 'configurations.html',
        controller: 'ConfigurationController'
    })
});
