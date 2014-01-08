var xrguide = angular.module('xrguide', ['ngRoute', 'xrguideControllers']);
xrguide.
    config(function ($routeProvider, $locationProvider) {
        'use strict';
        $routeProvider.
            when('/wares', {
                templateUrl: '/tmpl/wares.html',
                controller: 'WaresListCtrl'
            }).
            when('/', {
                templateUrl: '/tmpl/index.html'
            });
        $locationProvider.html5Mode(true);
    });

xrguide.service('Ware', ['$rootScope', function ($rootScope) {
    'use strict';
    var wares = [];
    var service = {
        wares: function() {
            return wares;
        },

        update: function ($http) {
            $http({
                method: 'GET',
                url: '/wares',
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                wares = data;
            });
            $rootScope.$broadcast('wares.update');
        }
    };

    return service;
}]);


xrguide.directive('wareRow', function () {
    'use strict';
    return {
        restrict: 'E',
        scope: {
            ware: '='
        },
        'templateUrl': '/tmpl/ware.html'
    };
});

var xrguideControllers = angular.module('xrguideControllers', []);

xrguideControllers.controller('WaresListCtrl', ['Ware', '$scope', '$http', function (Ware, $scope, $http) {
    'use strict';
    Ware.update($http);
    $scope.$watch('wares.update', function () {
        $scope.wares = Ware.wares();
    });
}]);

