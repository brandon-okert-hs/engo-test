package physics

import (
	"github.com/engoengine/math"

	"engo.io/engo"
	"github.com/bcokert/engo-test/logging"
)

func (e *ParticleEngine) Integrate(dt float32) {
	for _, p := range e.ParticleRegistry.particles {
		body := p.ParticleComponent()

		e.log.Debug("Before Integration", logging.F{"id": p.BasicEntity().ID(), "particleComponent": p.ParticleComponent()})

		// Calculate the net force on the object
		var netF engo.Point

		// Environment forces
		netF.Add(e.gravity)

		// Add other external forces on the object, then clear them
		netF.Add(body.ForceAccumulator)
		body.ForceAccumulator.Set(0, 0)

		// Acceleration from net force
		// a = F/m
		acceleration := netF
		acceleration.MultiplyScalar(body.InvMass)
		e.log.Debug("Integration Forces", logging.F{"gravity": e.gravity, "externalforces": body.ForceAccumulator, "netF": netF, "acceleration": acceleration})

		// Update the position
		// p = p + v*t + 0.5*a*t^2
		v := body.Velocity                  // v
		v.MultiplyScalar(dt)                // v*t
		body.SpaceComponent.Position.Add(v) // p = p + v*t
		a := acceleration                   // a
		a.MultiplyScalar(0.5 * dt * dt)     // 0.5*a*t^2
		body.SpaceComponent.Position.Add(a) // p = p + v*t + 0.5*a*t^2

		// Update the velocity
		// v = v*damp^t + a*t
		dampt := math.Pow(e.dampingFactor, dt) // damp^t
		body.Velocity.MultiplyScalar(dampt)    // v = v*damp^t
		a = acceleration                       // a
		a.MultiplyScalar(dt)                   // a*t
		body.Velocity.Add(a)                   // v = v*damp^t + a*t
		e.log.Debug("After Integration", logging.F{"id": p.BasicEntity().ID(), "particleComponent": p.ParticleComponent()})
	}
}
