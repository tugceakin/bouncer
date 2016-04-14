/**
 * Created by tugceakin on 3/13/16.
 */
//var bouncerApp = angular.module('bouncerApp', ['chart.js', 'ngRoute'])

var bouncerApp = angular.module('bouncerApp', ['chart.js', 'ngRoute', 'ngAnimate', 'uiSwitch'])

    .config(function ($routeProvider) {
    $routeProvider.
    when('/', {
        templateUrl: 'benchmarking.html',
        controller: 'RequestsGraphController'
    }).
    when('/benchmarking', {
        templateUrl: 'benchmarking.html',
        controller: 'RequestsGraphController'
    }).
    when('/page1', {
        templateUrl: 'page1.html',
        controller: 'ApiController'
    }).
    when('/page2', {
        templateUrl: 'page2.html',
        controller: 'ApiController'
    })
});
