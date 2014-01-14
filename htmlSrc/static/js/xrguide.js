var xrguide = angular.module('xrguide', ['ngRoute', 'xrguideControllers']);
xrguide.
    config(function ($routeProvider, $locationProvider) {
        'use strict';
        $routeProvider.
            when('/wares', {
                templateUrl: '/tmpl/wares.html',
                controller: 'WaresListCtrl'
            }).
            when('/ware/:wareId', {
                templateUrl: '/tmpl/ware.html',
                controller: 'WareDetailCtrl'
            }).
            when('/stations', {
                templateUrl: '/tmpl/stations.html',
                controller: 'StationsListCtrl'
            }).
            when('/', {
                templateUrl: '/tmpl/index.html'
            });
        $locationProvider.html5Mode(true);
    });

xrguide.service('XRGuide', ['$rootScope', '$http', function ($rootScope, $http) {
    'use strict';
    var wares = [],
        stations = {},
        service,
        getName,
        wareDefaults,
        stationDefaults;

    service = {
        getName: function (entity) {
            if (entity === undefined) {
                entity = this;
            }
            if (entity.Name !== undefined && entity.Name.Valid !== undefined && entity.Name.Valid) {
                return entity.Name.String;
            }
            if (entity.NameRaw !== undefined && entity.NameRaw.Valid !== undefined && entity.NameRaw.Valid) {
                return entity.NameRaw.String;
            }
            return '?';
        },

        wares: function () {
            return wares;
        },

        updateWares: function ($scope) {
            $http({
                method: 'GET',
                url: '/wares',
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                wares = angular.extend(data, wareDefaults);
                $scope.wares = wares;
                $rootScope.$broadcast('wares.update');
            });
        },

        getWare: function (wareId, $scope) {
            $http({
                method: 'GET',
                url: '/ware/' + wareId,
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                $scope.ware = angular.extend(data, wareDefaults);
                $scope.$broadcast('ware.update');
            });
        },

        updateStations: function ($scope) {
            $http({
                method: 'GET',
                url: '/stations',
                headers: {
                    'Accept': 'application/json'
                }
            }).success(function (data) {
                stations = angular.extend(data, stationDefaults);
                angular.forEach(stations, function (station) {
                    station.name = service.getName(station);
                });
                $scope.stations = stations;
                $rootScope.$broadcast('stations.update');
            });
        }
    };
    wareDefaults = {
        wareName: service.getName,

        wareSpecialist: function (ware) {
            if (ware === undefined) {
                ware = this;
            }
            if (ware.Specialist.Valid !== undefined && ware.Specialist.Valid) {
                return ware.Specialist.String;
            }
            return '-';
        },

        wareHasProduction: function (ware) {
            if (ware === undefined) {
                ware = this;
            }
            if (ware.Productions === undefined) {
                return false;
            }
            return true;
        }
    };
    stationDefaults = {
        stationName: service.getName
    };

    return service;
}]);

xrguide.directive('wareProduction', function () {
    return {
        restrict: 'E',
        templateUrl: '/tmpl/wareProduction.html'
    };
})

var xrguideControllers = angular.module('xrguideControllers', []);

xrguideControllers.controller('WaresListCtrl', ['XRGuide', '$scope', function (X, $scope) {
    'use strict';
    X.updateWares($scope);

    $scope.query = '';
    $scope.filterInternal = true;
    $scope.showTranspEquipment = false;
    $scope.showTranspInventory = false;
    $scope.showTranspEnergy = true;
    $scope.showTranspContainer = true;
    $scope.showTranspBulk = true;
    $scope.showTranspLiquid = true;
    $scope.showTranspFuel = true;

    $scope.isInternal = function (ware) {
        if (ware.Name.Valid !== undefined && ware.Name.Valid) {
            return false;
        }
        return true;
    };
    $scope.internalClass = function (ware) {
        if (!$scope.isInternal(ware)) {
            return '';
        }
        return 'ware-internal warning';
    };
    $scope.filterInternalExp = function (ware) {
        if (!$scope.filterInternal) {
            return true;
        }
        return !$scope.isInternal(ware);
    };
    $scope.filterTransports = function (ware) {
        if (!$scope.showTranspEquipment && ware.Transport === 'equipment') {
            return false;
        }
        if (!$scope.showTranspInventory && ware.Transport === 'inventory') {
            return false;
        }
        if (!$scope.showTranspEnergy && ware.Transport === 'energy') {
            return false;
        }
        if (!$scope.showTranspContainer && ware.Transport === 'container') {
            return false;
        }
        if (!$scope.showTranspBulk && ware.Transport === 'bulk') {
            return false;
        }
        if (!$scope.showTranspLiquid && ware.Transport === 'liquid') {
            return false;
        }
        if (!$scope.showTranspFuel && ware.Transport === 'fuel') {
            return false;
        }
        return true;
    };

//    $scope.$watch('wares.update', function () {});
}]);

xrguideControllers.controller('WareDetailCtrl', ['XRGuide', '$scope', '$routeParams', function (X, $scope, $routeParams) {
    'use strict';
    X.getWare($routeParams.wareId, $scope);
}]);

xrguideControllers.controller('StationsListCtrl', ['XRGuide', '$scope', function (X, $scope) {
    'use strict';
    X.updateStations($scope);
    $scope.order = 'name';
    $scope.getName = X.getName;
}]);