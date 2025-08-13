// SPDX-License-Identifier: EUPL-1.2

package routes

import (
	"errors"

	"git.zzdats.lv/edim/api-mdl/routes/responses"

	"azugo.io/azugo"
	"azugo.io/core/http"
	"github.com/valyala/fasthttp"
)

// @personId personID
// @title Get person data from CSDD
// @description Method return person driver licence data from CSDD
// @success 200 MDLResponse responses.MDLResponse "Get person data from CSDD"
// @failure 400 string string "Bad request"
// @failure 401 {empty} "Unauthorized"
// @failure 403 {empty} "Forbidden"
// @failure 404 {empty} "Not found"
// @failure 500 string string "Internal server error"
// @route /1.0/mdl [get].
func (r *router) mdl(ctx *azugo.Context) {
	csddresult, err := r.CsddService().GetCSDDData(ctx, ctx.User().Claim("code")[0])
	if err != nil {
		if errors.Is(err, http.NotFoundError{}) {
			ctx.StatusCode(fasthttp.StatusNotFound)
			ctx.Text("Data about drivers licence not found")

			return
		}

		ctx.Text(err.Error())
		ctx.Error(err)

		return
	}

	// skatamies vai ir atbildē "errors" bloks
	// ja ir, tad ir atbilde ar http 200, bet ar kļūdu
	if len(csddresult.Errors) > 0 {
		ctx.Error(errors.New(csddresult.Errors[0].ClientMessageCode + ": " + csddresult.Errors[0].ClientMessage))
		ctx.Text(csddresult.Errors[0].ClientMessageCode + ": " + csddresult.Errors[0].ClientMessage)
		ctx.StatusCode(fasthttp.StatusNotFound)

		return
	}

	mdlresult := &responses.MDLResponse{}
	mdlresult.PersonalAdministrativeNumber = ctx.User().Claim("code")[0]
	mdlresult.DocumentNumber = csddresult.Rowset[0].DocumentNumber
	mdlresult.BirthDate = csddresult.Rowset[0].BirthDate
	mdlresult.GivenName = csddresult.Rowset[0].GivenName
	mdlresult.FamilyName = csddresult.Rowset[0].FamilyName
	mdlresult.IssueDate = csddresult.Rowset[0].IssueDate
	mdlresult.ExpiryDate = csddresult.Rowset[0].ExpiryDate
	mdlresult.IssuingCountry = csddresult.Rowset[0].IssuingCountry
	mdlresult.IssuingAuthority = csddresult.Rowset[0].IssuingAuthority
	mdlresult.UnDistinguishingSign = csddresult.Rowset[0].UnDistinguishingSign
	mdlresult.Portrait = csddresult.Rowset[0].Portrait
	mdlresult.DrivingPrivileges = csddresult.Rowset[0].DrivingPrivileges

	ctx.JSON(mdlresult)
}
